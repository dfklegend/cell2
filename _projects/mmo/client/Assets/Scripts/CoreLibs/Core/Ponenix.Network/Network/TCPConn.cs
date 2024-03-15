using System;

namespace Phoenix.Network
{
    // 对消息的接收和发送
    public class TCPConnection
    {
        private long _id;
        public long id { get { return _id; } }

        TCPSession _session;
        public TCPSession session { get { return _session; } }        

        IProtocol _protocol;
        public IProtocol protocol { get { return _protocol; } }
        IChannel _channel;
        IMsgProcessor _processor;

        // Update触发
        Action<TCPConnection, IMsg> _cbProcessMsg;

        // 一些基于Connection的个性化对象
        private IConnectionImpl _impl;

        public TCPConnection(long id, TCPSession session,
            IChannel channel, IProtocol protocol, IMsgProcessor processor)
        {
            _id = id;
            _session = session;
            _channel = channel;
            _protocol = protocol;
            _processor = processor;
        }

        public T GetProtocol<T>()
            where T :class, IProtocol
        {
            return _protocol as T;
        }

        public void SetImpl(IConnectionImpl impl)
        {
            _impl = impl;
            _protocol.SetImpl(impl);
        }

        public T GetImpl<T>()
            where T: class, IConnectionImpl
        {
            if (_impl == null)
                return default(T);
            return _impl as T;
        }

        public bool IsReady()
        {
            return _protocol.IsReady();
        }

        public void SetCBProcessMsg(Action<TCPConnection, IMsg> cb)
        {
            _cbProcessMsg = cb;
        }

        public void Start()
        {            
            _session.Start();
            _protocol.OnConnected();
        }

        public void Stop()
        {
            if(_session == null)
            {
                Env.L.Error("Stop already");
                return;
            }
            _session.Stop();
            _session = null;
            _protocol.Stop();
            _protocol = null;
            _channel = null;
            _processor = null;
            _cbProcessMsg = null;
        }

        public void Update()
        {            
            _protocol.Update();
            // protcol里面可能触发异步函数回调
            // 然后可能导致上面的Stop
            if (_protocol == null)
              return;            
            processMsgs();
            _channel.Update();
        }

        private void processMsgs()
        {
            IMsg msg;
            while (_processor != null)
            {
                msg = _processor.Pop();
                if (msg == null)
                    break;                
                _cbProcessMsg?.Invoke(this, msg);
            }
        }

        private void testMsg(IMsg msg)
        {
            SimpleMsg simple = msg as SimpleMsg;
            if (simple == null)
                return;
            //Console.WriteLine($"Got: {simple.content}");
            if (simple.content == "Ack")
                return;

            // send back
            var ack = new SimpleMsg();
            ack.content = "Ack";
            _protocol.SendMsg(ack);
        }
    }
}

