namespace Phoenix.Network.Protocol.Pomelo
{
    // 协议转发给Protocol
    // API的Context
    public class PomeloSession : IClientSession, API.IContext
    {
        private BasePomeloProtocol _protocol;

        public PomeloSession(BasePomeloProtocol protocol)
        {
            _protocol = protocol;
        }

        public int GetHandle()
        {
            return _protocol.GetHandle();
        }

        public PomeloNode node { get { return _protocol.node; } }

        public T GetImpl<T>()
            where T : class, IConnectionImpl
        {
            return _protocol.GetImpl<T>();
        }

        public BasePomeloProtocol GetProtocol()
        {
            return _protocol;
        }

        public bool IsConnected()
        {
            return _protocol.IsReady();
        }

        public void SendNotify(string route, byte[] data)
        {
            _protocol.sendNotify(route, data);
        }

        public void SendResponse(uint reqId, byte[] data)
        {
            _protocol.sendResponse(reqId, data);
        }

        public void SendRequest(string route, uint reqId, byte[] data)
        {
            _protocol.sendRequest(route, reqId, data);
        }
    }
}

