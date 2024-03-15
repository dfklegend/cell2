using System;
using System.Collections.Generic;
using System.Linq;
using System.Net.Sockets;
using System.Threading;
using Phoenix.Scheduler;

namespace Phoenix.Network
{
    // 管理Acceptor的connection
    public class TCPConnections
    {
        private ThreadSynchronizationContext _main;
        
        private Dictionary<long, TCPConnection> _conns = new Dictionary<long, TCPConnection>();
        // 方便遍历
        private List<TCPConnection> _quickList = new List<TCPConnection>();
        private bool _dirtyConns = false;

        private NetConfig _config;

        public Action<TCPConnection> cbOnSessionRemove;

        public TCPConnections(ThreadSynchronizationContext context, NetConfig config)
        {
            _config = config;
            _main = context;
        }

        public void AddConnection(TCPConnection conn)
        {
            Env.L.Info($"AddConnection {conn.id} Thread:{Thread.CurrentThread.ManagedThreadId}");
            _conns[conn.id] = conn;
            _dirtyConns = true;
        }

        public void RemoveConnection(TCPConnection conn)
        {
            long id = conn.id;

            cbOnSessionRemove?.Invoke(conn);
            conn.Stop();
            _conns.Remove(id);
            _dirtyConns = true;

            Env.L.Info($"RemoveSession {id} Thread:{Thread.CurrentThread.ManagedThreadId}");
        }

        public void Update()
        {
            if(_dirtyConns)
            {
                makeQuickList();
                _dirtyConns = false;
            }            
            
            for(var i = 0; i < _quickList.Count; i ++)
            {
                _quickList[i].Update();
            }
        }

        public void Stop()
        {
            foreach(var conn in _conns.Values)
            {
                conn.Stop();
            }
            _conns.Clear();
            _dirtyConns = true;
        }

        private void makeQuickList()
        {
            _quickList = _conns.Values.ToList();
        }

        private long allocId()
        {
            return ConnIDAllocer.Alloc();
        }

        public TCPConnection onAccept(Socket socket)
        {
            long id = allocId();
            var conn = TCPBuilder.Build(this._config, id, socket);
            AddConnection(conn);

            conn.session.SetCBError(onSessionErr);
            conn.Start();
            return conn;
        }
        
        public TCPConnection GetConnection(long id)
        {
            TCPConnection session;
            if (!_conns.TryGetValue(id, out session))
                return null;
            return session;
        }

        private void onSessionErr(long id, SocketError err)
        {
            Env.L.Info($"onSessionErr {id} Thread:{Thread.CurrentThread.ManagedThreadId}");

            _main.PostAlways((state) => {
                var session = GetConnection(id);
                if (session != null)
                    RemoveConnection(session);
            }, null);
        }
    }
}

