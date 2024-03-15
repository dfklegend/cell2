using System;
using System.IO;

namespace Phoenix.Network
{
    // 4字节长度
    // 长度字节
    public class SimpleMsg : IMsg
    {
        public const int MAX_MSG_LEN = 2 * 1024 * 1024;
        public const int LENGTH_SIZE = sizeof(UInt32);

        
        public UInt32 length;
        public string content;

        public void Read(byte[] bytes)
        {
            MsgStream ms = new MsgStream(bytes);
            Deserialize(ms);
        }

        public void Serialize(MsgStream ms)
        {
            ms.Write(length);
            ms.DWriteUTF8String32(content);

            length = (UInt32)ms.Length;
            ms.Seek(0, SeekOrigin.Begin);
            ms.Write(length);
            ms.Seek(0, SeekOrigin.Begin);
        }

        public void Deserialize(MsgStream ms)
        {
            length = ms.ReadUInt32();
            content = ms.DReadUTF8String32();
        }
    }

    public class SimpleMsgDecoder : IMsgDecoder
    {
        private UInt32 length = 0;
        public byte[] _lengthBytes = new byte[SimpleMsg.LENGTH_SIZE];

        public IMsg Make(Stream stream)
        {
            if(length > 0)
            {
                return continueMake(stream);
            }
            if (stream.Length < SimpleMsg.LENGTH_SIZE)
                return null;
            stream.Read(_lengthBytes, 0, SimpleMsg.LENGTH_SIZE);
            length = BitConverter.ToUInt32(_lengthBytes, 0);
            
            if(length > SimpleMsg.MAX_MSG_LEN)
            {
                Env.L.Error("SimpleMsgDecoder Error msg, skip!");
                length = 0;
                return null;
            }
            return continueMake(stream);
        }

        public IMsg continueMake(Stream stream)
        {
            // 暂时还不够
            if (stream.Length < length - SimpleMsg.LENGTH_SIZE)
                return null;

            byte[] byMsg = new byte[length];
            Array.Copy(_lengthBytes, byMsg, SimpleMsg.LENGTH_SIZE);
            stream.Read(byMsg, SimpleMsg.LENGTH_SIZE, (int)length - SimpleMsg.LENGTH_SIZE);
            
            prepareNextMsg();


            SimpleMsg msg = new SimpleMsg();
            msg.Read(byMsg);
            return msg;
        }

        private void prepareNextMsg()
        {
            this.length = 0;
        }
    }

    public class SimpleMsgEncoder : IMsgEncoder
    {
        private static SimpleMsgEncoder _it;
        public static SimpleMsgEncoder It 
        { 
            get 
            {
                if (_it == null)
                    _it = new SimpleMsgEncoder();
                return _it;
            } 
        }

        public Stream Write(IMsg msg)
        {
            SimpleMsg simple = msg as SimpleMsg;
            if (simple == null)
                return null;
            MsgStream ms = new MsgStream();           
            simple.Serialize(ms);
            return ms;
        }
    }

    public class SimpleCoderFactory : IMsgCoderFactory
    {
        public IMsgDecoder CreateDecoder()
        {
            return new SimpleMsgDecoder();
        }

        public IMsgEncoder CreateEncoder()
        {
            return new SimpleMsgEncoder();
        }
    }

}

