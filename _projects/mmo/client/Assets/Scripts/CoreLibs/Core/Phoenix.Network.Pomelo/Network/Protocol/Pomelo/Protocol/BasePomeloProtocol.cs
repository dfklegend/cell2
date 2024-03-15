using System;
using System.IO;
using Pomelo.DotNetClient;
using SimpleJson;
using TimeUtil = Phoenix.Utils.TimeUtil;

namespace Phoenix.Network.Protocol.Pomelo
{
    /*
     * c->s handshake
     * s->c handshake
     * c->s handshakeack
     */
    public class BasePomeloProtocol : IProtocol
    {
        private IChannel _channel;
        private IMsgEncoder _encoder;
        private IMsgDecoder _decoder;
        private IMsgProcessor _processor;
        protected MessageProtocol _messageProtocol;

        protected PomeloNode _node;
        QueueProcessor _queue = new QueueProcessor();

        protected int _heartBeatInterval = 0;
       

        private bool _stopped = false;

        IConnectionImpl _impl;
        
        public Action cbOnReady;

        public void Init(IChannel channel,
            IMsgEncoder encoder, IMsgDecoder decoder,
            IMsgProcessor processor, IProtocolConfig config)
        {
            _channel = channel;
            _encoder = encoder;
            _decoder = decoder;
            _processor = processor;

            _node = new PomeloNode();
            _node.Init(new PomeloSession(this), new MsgSender());
            applyConfig(config as PomeloConfig);

            {
                // debug
                MsgDecoder pd = decoder as MsgDecoder;
                pd.debugSocketHandle = _channel.GetHandle();
            }
        }

        public int GetHandle()
        {
            return _channel.GetHandle();
        }

        private void applyConfig(PomeloConfig config)
        {
            if (config == null)
                return;            
            if (config.dispatcher != null)
                _node.SetMsgDispatcher(config.dispatcher);           
            _node.sender.SetSerializer(config.serializer);
        }

        public PomeloNode node { get { return _node; } }

        public void MakeMsg(Stream stream)
        {
            IMsgDecoder decoder = _decoder;
            while (true)
            {
                IMsg msg;

                long beforeLen = stream.Length;
                msg = decoder.Make(stream);
                if (msg == null)
                {
                    //Env.L.FileLog($"{GetHandle()} makeMsg null, consume {beforeLen - stream.Length}");
                    break;
                }

                //Env.L.FileLog($"{GetHandle()} makeMsg succ, consume {beforeLen- stream.Length}");
                if (processInternalMsg(msg))
                    continue;
                // decode message
                PomeloMsg pomeloMsg = msg as PomeloMsg;
                if (pomeloMsg.package.type == PackageType.PKG_DATA)
                {
                    if(_messageProtocol == null)
                    {
                        Env.L.Error("Error! messageProtocol not init");
                        //Env.L.FileLog("Error! messageProtocol not init");
                        break;
                    }
                    pomeloMsg.message = _messageProtocol.decode(pomeloMsg.package.body);
                }
                
                //Env.L.FileLog($"{GetHandle()} got type: {pomeloMsg.message.type} reqId: {pomeloMsg.message.id}");
                                
                _queue.Process(this, msg);
            }
        }

        public void SendMsg(IMsg msg)
        {            
        }

        public virtual void OnConnected()
        {
            _node.Start();
        }

        public virtual bool IsReady()
        {
            return true;
        }

        public void SetCBOnReady(Action cb)
        {
            cbOnReady = cb;
        }

        public void OnReady()
        {
            cbOnReady?.Invoke();
        }

        // 目前在线程内执行
        private bool processInternalMsg(IMsg m)
        {
            PomeloMsg msg = m as PomeloMsg;
            var pkt = msg.package;

            //Console.WriteLine($"got data: {pkt.type}");
            switch (pkt.type)
            {
                case PackageType.PKG_HANDSHAKE:
                    processHandshake(msg);
                    return true;
                case PackageType.PKG_HANDSHAKE_ACK:
                    processHandshakeAck(msg);
                    return true;
                case PackageType.PKG_HEARTBEAT:
                    processHeartbeat(msg);
                    return true;
                case PackageType.PKG_KICK:
                    processKick(msg);
                    return true;
                case PackageType.PKG_DATA:
                    // do next
                    //Console.WriteLine("got data");
                    return false;
            }            
            return false;
        }

        public virtual void Update()
        {
            if (_stopped)
                return;            
            processMsgs();
            _node.Update();
        }

        protected void processMsgs()
        {
            IMsg msg;
            while (!_stopped)
            {
                msg = _queue.Pop();
                if (msg == null)
                    break;
                PomeloMsg pomeloMsg = msg as PomeloMsg;
                _node.Process(pomeloMsg.message);
            }
        }        

        public void Stop() 
        {
            _stopped = true;
            _node.Stop();
        }

        public bool IsStopped()
        {
            return _stopped;
        }

        public void ReqStopSession()
        {
            if (_channel != null)
            {
                _channel.StopSession();
            }
        }

        protected virtual void processHandshake(PomeloMsg msg)
        {            
        }

        protected virtual void processHandshakeAck(PomeloMsg msg)
        {
        }

        protected virtual void processHeartbeat(PomeloMsg msg)
        {           
        }

        protected virtual void processKick(PomeloMsg msg)
        {

        }
        

        //Send request, user request id 
        internal void send(string route, uint id, JsonObject msg)
        {
            if (msg == null)
                msg = new JsonObject();
            byte[] body = _messageProtocol.encode(route, id, msg);

            send(PackageType.PKG_DATA, body);
        }

        internal void send(PackageType type)
        {   
            _channel.SendMsg(PackageProtocol.encode(type));
        }

        internal void send(PackageType type, byte[] body)
        {
            byte[] pkg = PackageProtocol.encode(type, body);
            transporter_send(pkg);
        }

        void transporter_send(byte[] pkg)
        {
            _channel.SendMsg(pkg);
        }

        internal void sendRequest(string route, uint id, byte[] body)
        {            
            byte[] msgBody = _messageProtocol.newEncode(MessageType.MSG_REQUEST, route, false, id, body);
            send(PackageType.PKG_DATA, msgBody);
        }

        internal void sendNotify(string route, byte[] body)
        {
            byte[] msgBody = _messageProtocol.newEncode(MessageType.MSG_NOTIFY, route, false, 0, body);
            send(PackageType.PKG_DATA, msgBody);
        }

        internal void sendResponse(uint id, byte[] body)
        {
            byte[] msgBody = _messageProtocol.newEncode(MessageType.MSG_RESPONSE, "", false, id, body);
            send(PackageType.PKG_DATA, msgBody);
        }

        public void SetImpl(IConnectionImpl impl) 
        { 
            _impl = impl; 
        }

        public T GetImpl<T>() 
            where T : class, IConnectionImpl 
        {
            if (_impl == null)
                return default(T);
            return _impl as T;
        }
    }
}

