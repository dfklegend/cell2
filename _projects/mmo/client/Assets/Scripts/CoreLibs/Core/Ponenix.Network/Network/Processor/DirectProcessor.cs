using System;

namespace Phoenix.Network
{
    // 缓存，等外部处理
    public class DirectProcessor : IMsgProcessor
    {
        private Action<IProtocol, IMsg> _cbProcess;

        public void SetCBProcess(Action<IProtocol, IMsg> cb)
        {
            _cbProcess = cb;
        }

        public void Process(IProtocol protocol, IMsg msg)
        {
            _cbProcess(protocol, msg);
        }

        public void Update()
        {
        }

        public IMsg Pop()
        {
            return null;
        }
    }

    public class DirectProcessorFactory : IProcessorFactory
    {
        public IMsgProcessor Create()
        {
            return new DirectProcessor();
        }
    }
}

