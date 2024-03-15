namespace Phoenix.Network
{
    // 缓存，等外部处理
    public class EmptyProcessor : IMsgProcessor
    {
        public void Process(IProtocol protocol, IMsg msg)
        { 
        }

        public void Update()
        {
        }

        public IMsg Pop()
        {
            return null;
        }
    }

    public class EmptyProcessorFactory : IProcessorFactory
    {
        public IMsgProcessor Create()
        {
            return new EmptyProcessor();
        }
    }
}

