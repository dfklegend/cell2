using System;
using System.IO;
using Pomelo.DotNetClient;

namespace Phoenix.Network.Protocol.Pomelo
{
    /**
    * Package protocol encode.
    *
    * Pomelo package format:
    * +------+-------------+------------------+
    * | type | body length |       body       |
    * +------+-------------+------------------+
    *
    * Head: 4bytes
    *   0: package type,
    *      1 - handshake,
    *      2 - handshake ack,
    *      3 - heartbeat,
    *      4 - data
    *      5 - kick
    *   1 - 3: big-endian body length
    * Body: body length bytes
    *
    * @param  {Number}    type   package type
    * @param  {ByteArray} body   body content in bytes
    * @return {ByteArray}        new byte array that contains encode result
    */

    public class PomeloMsg : IMsg
    {
        public Package package;
        public Message message;
        public PomeloMsg(Package p)
        {
            package = p;
        }
    }

    public class MsgDecoder : IMsgDecoder
    {
        private int packetType = 0;
        private UInt32 length = 0;
        public byte[] _lengthBytes = new byte[PomeloDefine.LENGTH_SIZE];
        public byte[] _tempBytes = new byte[PomeloDefine.LENGTH_SIZE];


        public int debugSocketHandle = 0;

        public IMsg Make(Stream stream)
        {
            if (length > 0)
            {
                return continueMake(stream);
            }
            if (stream.Length < SimpleMsg.LENGTH_SIZE)
                return null;
            stream.Read(_lengthBytes, 0, SimpleMsg.LENGTH_SIZE);
            
            UInt32 data = BitConverter.ToUInt32(_lengthBytes, 0);
            
            packetType = (int)(data & 0xFF);

            _tempBytes[0] = (byte)((data >> 24) & 0xFF);
            _tempBytes[1] = (byte)((data >> 16) & 0xFF);
            _tempBytes[2] = (byte)((data >> 8) & 0xFF);
            _tempBytes[3] = 0;
            
            // bodyLength + 4
            this.length = BitConverter.ToUInt32(_tempBytes, 0) + PomeloDefine.LENGTH_SIZE;


            if (length > SimpleMsg.MAX_MSG_LEN)
            {
                Env.L.Error("Pomelo MsgDecoder too large msg length, skip!");
                //Env.L.FileLog($"{debugSocketHandle} Error msg, skip!");
                length = 0;
                return null;
            }
            return continueMake(stream);
        }

        public IMsg continueMake(Stream stream)
        {
            // 暂时还不够
            if (stream.Length < length - PomeloDefine.LENGTH_SIZE)
                return null;

            byte[] byMsg = new byte[length];            
            Array.Copy(_lengthBytes, byMsg, SimpleMsg.LENGTH_SIZE);
            // body 
            stream.Read(byMsg, PomeloDefine.LENGTH_SIZE, (int)length - PomeloDefine.LENGTH_SIZE);

            Package pkg = PackageProtocol.decode(byMsg);

            prepareNextMsg();
            
            return new PomeloMsg(pkg);
        }

        private void prepareNextMsg()
        {
            this.length = 0;
        }
    }    

    public class PomeloCoderFactory : IMsgCoderFactory
    {
        public IMsgDecoder CreateDecoder()
        {
            return new MsgDecoder();
        }

        public IMsgEncoder CreateEncoder()
        {
            return null;
        }
    }        
}

