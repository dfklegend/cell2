using System;
using System.Collections.Generic;
using System.Text;
using SimpleJson;
namespace Pomelo.DotNetClient
{
    public class MessageProtocol
    {
        private Dictionary<string, ushort> dict = new Dictionary<string, ushort>();
        private Dictionary<ushort, string> abbrs = new Dictionary<ushort, string>();
        private JsonObject encodeProtos = new JsonObject();
        private JsonObject decodeProtos = new JsonObject();
        private Dictionary<uint, string> reqMap;
        private Protobuf.Protobuf protobuf;

        public const int MSG_Route_Limit = 255;
        public const int MSG_Compress_Route_Mask = 0x01;
        public const int MSG_Type_Mask = 0x07;
        public const int MSG_Err_Mask = 0x20;

        public MessageProtocol(JsonObject dict, JsonObject serverProtos, JsonObject clientProtos)
        {
            ICollection<string> keys = dict.Keys;

            foreach (string key in keys)
            {
                ushort value = Convert.ToUInt16(dict[key]);
                this.dict[key] = value;
                this.abbrs[value] = key;
            }

            protobuf = new Protobuf.Protobuf(clientProtos, serverProtos);
            this.encodeProtos = clientProtos;
            this.decodeProtos = serverProtos;

            this.reqMap = new Dictionary<uint, string>();
        }

        public byte[] encode(string route, JsonObject msg)
        {
            return encode(route, 0, msg);
        }

        public byte[] encode(string route, uint id, JsonObject msg)
        {
            return doEncode(route, id, -1, msg, null);
        }

        public byte[] encode(string route, uint id, byte[] msgData)
        {
            return doEncode(route, id, -1, null, msgData);
        }

        public byte[] doEncode(string route, uint id, int messageType, JsonObject msg, byte[] bodyData)
        {
            if(messageType == -1)
            {
                if (id > 0)
                    messageType = (int)MessageType.MSG_REQUEST;
                else
                    messageType = (int)MessageType.MSG_NOTIFY;
            }

            if(msg != null)
            {
                if (encodeProtos.ContainsKey(route))
                {
                    bodyData = protobuf.encode(route, msg);
                }
                else
                {
                    bodyData = Encoding.UTF8.GetBytes(msg.ToString());
                }                
            }

            return newEncode((MessageType)messageType, route, false, id, bodyData);
        }        

        public Message decode(byte[] buffer)
        {
            return doDecode(buffer, true);
        }       

        public Message doDecode(byte[] buffer, bool needDeserialize)
        {
            var msg = newDecode(buffer);
            // 暂时去掉此处的反序列化
            if (false && needDeserialize)
            {
                if (decodeProtos.ContainsKey(msg.route))
                {
                    msg.oldJsonData = protobuf.decode(msg.route, msg.rawData);
                }
                else
                {
                    if (msg.rawData.Length > 0)
                    {
                        var content = Encoding.UTF8.GetString(msg.rawData);                        
                        msg.oldJsonData = Phoenix.Network.Protocol.Pomelo.PomeloUtil.DeserializeObject(content);
                    }
                }                
            }
            return msg;
        }

        private void writeInt(int offset, uint value, byte[] bytes)
        {
            bytes[offset] = (byte)(value >> 24 & 0xff);
            bytes[offset + 1] = (byte)(value >> 16 & 0xff);
            bytes[offset + 2] = (byte)(value >> 8 & 0xff);
            bytes[offset + 3] = (byte)(value & 0xff);
        }

        private void writeShort(int offset, ushort value, byte[] bytes)
        {
            bytes[offset] = (byte)(value >> 8 & 0xff);
            bytes[offset + 1] = (byte)(value & 0xff);
        }

        private ushort readShort(int offset, byte[] bytes)
        {
            ushort result = 0;

            result += (ushort)(bytes[offset] << 8);
            result += (ushort)(bytes[offset + 1]);

            return result;
        }

        private int byteLength(string msg)
        {
            return Encoding.UTF8.GetBytes(msg).Length;
        }

        private void writeBytes(byte[] source, int offset, byte[] target)
        {
            for (int i = 0; i < source.Length; i++)
            {
                target[offset + i] = source[i];
            }
        }

        bool msgHasId(MessageType type)
        {
            switch(type)
            {
                case MessageType.MSG_REQUEST:
                case MessageType.MSG_RESPONSE:
                    return true;
            }
            return false;
        }

        bool msgHasRoute(MessageType type)
        {
            switch(type)
            {
                case MessageType.MSG_REQUEST:
                case MessageType.MSG_NOTIFY:
                case MessageType.MSG_PUSH:
                    return true;
            }
            return false;
        }

        public Message newDecode(byte[] buffer)
        {
            // Decode head
            //Get flag
            byte flag = buffer[0];
            //Set offset to 1, for the 1st byte will always be the flag
            int offset = 1;

            // flag 一个byte
            // MMMMTTTR  高4位MMMM 有err和gzip压缩标志 TTT是type(request,notify,...), R是route压缩标志
            //Get type from flag;
            MessageType type = (MessageType)((flag >> 1) & MSG_Type_Mask);
            bool compressRoute = (flag & MSG_Compress_Route_Mask) != 0;
            bool err = (flag & MSG_Err_Mask) != 0;

            uint id = 0;
            string route = "";

            if (msgHasId(type))
            {
                int length;
                id = (uint)Protobuf.Decoder.decodeUInt32(offset, buffer, out length);

                offset += length;
            }
           
            if(msgHasRoute(type))
            {
                if(compressRoute)
                {
                    ushort routeId = readShort(offset, buffer);

                    // TODO: safety
                    route = abbrs[routeId];

                    offset += 2;
                }
                else 
                {
                    byte length = buffer[offset];
                    offset += 1;

                    route = Encoding.UTF8.GetString(buffer, offset, length);
                    offset += length;
                }
            }

            byte[] body = new byte[buffer.Length - offset];
            for (int i = 0; i < body.Length; i++)
            {
                body[i] = buffer[i + offset];
            }

            return new Message(type, id, route, body, err);
        }


        public byte[] newEncode(MessageType type, string route, bool compressRoute, uint id, byte[] bodyData)        
        {
            int routeLength = byteLength(route);
            if (routeLength > MSG_Route_Limit)
            {
                throw new Exception("Route is too long!");
            }

            //Encode head
            //The maximus length of head is 1 byte flag + 5 bytes message id + route string length + 1byte
            // 动态id的字节最大5字节
            const int MAX_BONUS_LENGTH = 7;
            byte[] head = new byte[routeLength + MAX_BONUS_LENGTH];
            int offset = 1;
            byte flag = 0;

            if (msgHasId(type))
            {
                byte[] bytes = Protobuf.Encoder.encodeUInt32(id);
                writeBytes(bytes, offset, head);
                offset += bytes.Length;
            }

            if(msgHasRoute(type))
            {
                // 压缩route
                if(compressRoute && dict.ContainsKey(route))
                {
                    ushort cmpRoute = dict[route];
                    writeShort(offset, cmpRoute, head);
                    flag |= MSG_Compress_Route_Mask;
                    offset += 2;
                }
                else
                {
                    //Write route length
                    head[offset++] = (byte)routeLength;

                    //Write route
                    writeBytes(Encoding.UTF8.GetBytes(route), offset, head);
                    offset += routeLength;
                }
            }

            flag |= (byte)(((byte)type) << 1);
            head[0] = flag;

            //Encode body
            // protobuf.encode和Json的序列化重新包装接口处理
            byte[] body = bodyData;

            //Construct the result
            byte[] result = new byte[offset + body.Length];
            for (int i = 0; i < offset; i++)
            {
                result[i] = head[i];
            }

            for (int i = 0; i < body.Length; i++)
            {
                result[offset + i] = body[i];
            }
            return result;
        }
    }
}