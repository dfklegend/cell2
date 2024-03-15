using System;
using System.Net;
using System.Net.Sockets;
using Phoenix.Network.Protocol.Pomelo;
using Phoenix.Scheduler;

namespace Phoenix.Network
{
    // 便于使用
    public class PomeloClient
    {
        private RichMsgDispatcher _dispatcher = new RichMsgDispatcher();
        public RichMsgDispatcher dispatcher { get { return _dispatcher; } }
        private TCPClient _client;
        private int _runId = -1;
        private PomeloNode _node;
        public PomeloNode node { get{return _node;} }

        // 可添加
        public Action<SocketError> cbConnectFail { set { _client.cbConnectFail = value; } }
        public Action<SocketError> cbConnectionErr { set { _client.cbConnectionErr = value; } }
        public Action<PomeloClient> cbConnected;

        public PomeloClient(string apiCategories, Serializer.ISerializer serializer = null)
        {
            if(serializer == null)
                serializer = new Serializer.SimpleJsonSerializer();
            _dispatcher.SetSerializer(serializer);
            _dispatcher.CreateServiceImpl(apiCategories);

            var config = new NetConfig();
            initPomeloClientConfig(config, serializer, _dispatcher);
            _client = new TCPClient(ThreadMgr.it.mainContext, config);
            
            _client.cbReady = (client) => {
                onConnectted();
            };

            _runId = ThreadMgr.it.mainThread.AddRunStep(()=> { onUpdate(); });
        }

        static void initPomeloClientConfig(NetConfig config, 
            Serializer.ISerializer serializer, RichMsgDispatcher dispatcher)
        {
            var pconfig = new PomeloConfig();
            pconfig.WithMsgDispatcher(dispatcher);
            pconfig.WithSerialize(serializer);            

            config.SetCoderFactory(new PomeloCoderFactory())
                .SetProtocolFactory(new ClientProtocolFactory())
                .WithProtocolConfig(pconfig)
                .SetProcessorFactory(new EmptyProcessorFactory());
        }

        // 开始连接
        public void Connect(string ip, int port)
        {
            var ipEndPoint = new IPEndPoint(IPAddress.Parse(ip), port);
            Connect(ipEndPoint);
        }

        public void Connect(IPEndPoint ipEndPoint)
        {            
            _client.Connect(ipEndPoint);
        }        

        public bool IsConnected()
        {
            return _client.IsConnected();
        }

        public void ReConnect()
        {
            _client.ReConnect();
        }

        private void onConnectted()
        {
            _node = _client.GetProtocol<BasePomeloProtocol>().node;
            cbConnected?.Invoke(this);
        }

        private void onUpdate()
        {
            _client.Update();
        }

        public void Close()
        {
            if (_runId > 0)
            {
                ThreadMgr.it.mainThread.RemoveRunStep(_runId);
                _runId = -1;
            }
            _client.Stop();
        }
    }
}
