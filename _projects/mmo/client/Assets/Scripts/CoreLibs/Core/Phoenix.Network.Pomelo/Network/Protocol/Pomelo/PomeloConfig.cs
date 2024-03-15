namespace Phoenix.Network.Protocol.Pomelo
{
    // 协议额外的初始化
    public class PomeloConfig : IProtocolConfig
    {
        public IMsgDispatcher dispatcher;        
        public Serializer.ISerializer serializer;
        public string APICategoris = "";

        public PomeloConfig WithMsgDispatcher(IMsgDispatcher d)
        {
            dispatcher = d;
            return this;
        }

        public PomeloConfig WithSerialize(Serializer.ISerializer s)
        {
            serializer = s;
            return this;
        }

        public PomeloConfig WithAPICategories(string categoris)
        {
            APICategoris = categoris;
            return this;
        }
    }
}

