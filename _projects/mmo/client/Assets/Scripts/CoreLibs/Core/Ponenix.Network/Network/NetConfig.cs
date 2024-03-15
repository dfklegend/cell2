using System;

namespace Phoenix.Network
{
    // 便于构建
    public class NetConfig
    {
        public IMsgCoderFactory msgCoderFactory;
        public IProtocolFactory protocolFactory;
        public IProtocolConfig protocolConfig;
        public IProcessorFactory processorFactory;
        public Action<IMsgProcessor> cbInitProcessor;
        

        public NetConfig SetCoderFactory(IMsgCoderFactory factory)
        {
            msgCoderFactory = factory;
            return this;
        }

        public NetConfig SetProtocolFactory(IProtocolFactory factory)
        {
            protocolFactory = factory;
            return this;
        }

        public NetConfig WithProtocolConfig(IProtocolConfig config)
        {
            protocolConfig = config;
            return this;
        }

        public NetConfig SetProcessorFactory(IProcessorFactory factory)
        {
            processorFactory = factory;
            return this;
        }

        public NetConfig WithProcessorInitor(Action<IMsgProcessor> init)
        {
            cbInitProcessor = init;
            return this;
        }
    }    
}

