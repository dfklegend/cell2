using System;
using System.IO;

namespace Phoenix.Network
{
    public class SimpleProtocol : IProtocol
    {
        private IChannel _channel;
        private IMsgEncoder _encoder;
        private IMsgDecoder _decoder;
        private IMsgProcessor _processor;

        public void Init(IChannel channel,
            IMsgEncoder encoder, IMsgDecoder decoder,
            IMsgProcessor processor, IProtocolConfig config)
        {
            _channel = channel;
            _encoder = encoder;
            _decoder = decoder;
            _processor = processor;
        }

        public int GetHandle()
        {
            return _channel.GetHandle();
        }

        public void MakeMsg(Stream stream)
        {
            IMsgDecoder decoder = _decoder;
            while (true)
            {
                IMsg msg;
                msg = decoder.Make(stream);
                if (msg == null)
                    break;

                _processor.Process(this, msg);
            }
        }

        public void SendMsg(IMsg msg)
        {
            _channel.SendMsg(_encoder.Write(msg));
        }

        public void OnConnected()
        {
        }

        public bool IsReady()
        {
            return true;
        }

        public void SetCBOnReady(Action cb)
        {

        }

        public void Update() { }

        public void Stop() { }

        public void SetImpl(IConnectionImpl impl) { }
        public T GetImpl<T>() where T : class, IConnectionImpl { return null; }
    }

    public class SimpleProtocolFactory : IProtocolFactory
    {
        public IProtocol Create()
        {
            return new SimpleProtocol();
        }
    }
}

