using Pomelo.DotNetClient;

namespace Phoenix.Network.Protocol.Pomelo
{
    public interface IMsgDispatcher
    {
        void OnSessionStart(IClientSession session);
        void OnSessionStop(IClientSession session);
        void Process(IClientSession session, Message msgIn);
        Serializer.ISerializer GetSerializer();
    }    
}

