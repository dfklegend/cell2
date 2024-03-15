using System.Net.Sockets;

namespace Phoenix.Network
{
    // 对消息的接收和发送
    public static class TCPBuilder
    {
        public static TCPConnection Build(NetConfig config, long id, Socket socket)
        {
            var session = new TCPSession(id, socket, null);
            var protocol = config.protocolFactory.Create();
            var processor = config.processorFactory.Create();
            var channel = new Channel(session, protocol);
            
            protocol.Init(channel, config.msgCoderFactory.CreateEncoder(),
                config.msgCoderFactory.CreateDecoder(), processor, config.protocolConfig);
            config.cbInitProcessor?.Invoke(processor);
            
            session.SetChannel(channel);
            var conn = new TCPConnection(id, session, channel, protocol, processor);            
            return conn;
        }
    }
}

