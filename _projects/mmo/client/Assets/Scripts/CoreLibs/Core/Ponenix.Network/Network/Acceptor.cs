using System;
using System.Net;
using System.Net.Sockets;
using System.Threading;
using Phoenix.Scheduler;

namespace Phoenix.Network
{
    public class Acceptor
    {
        private Socket _socket;
        private readonly SocketAsyncEventArgs innArgs = new SocketAsyncEventArgs();
        
        private ThreadSynchronizationContext _main;
        private Action<Socket> _onAccept;

        public bool Start(IPEndPoint ipEndPoint, ThreadSynchronizationContext main,
            Action<Socket> cbAccept)
        {
            this._main = main;

            if (!doListen(ipEndPoint))
            {   
                Env.L.Error($"listen error: {ipEndPoint}");
                return false;
            }
            this.innArgs.Completed += this.OnComplete;
            _onAccept = cbAccept;

            // 保证OnAcceptComplete在主线程执行
            _main.Post((state) =>
            {
                acceptAsync();
            }, null);
            return true;
        }

        private bool doListen(IPEndPoint ipEndPoint)
        {
            try 
            {
                this._socket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
                this._socket.Bind(ipEndPoint);
                this._socket.Listen(128);
                return true;
            }
            catch(Exception e)
            {                
                Phoenix.Utils.SystemUtil.LogHandledException(e);
                return false;
            }
        }

        public void Stop()
        {
            if (this._socket == null)
                return;
            this._socket.Close();
            this._socket = null;
            this.innArgs.Dispose();
            this._main = null;
            
        }

        private void OnComplete(object sender, SocketAsyncEventArgs e)
        {
            switch (e.LastOperation)
            {
                case SocketAsyncOperation.Accept:
                    SocketError socketError = e.SocketError;
                    Socket acceptSocket = e.AcceptSocket;

                    // 成功accept
                    _main?.Post((state) => 
                    {
                        this.OnAcceptComplete(socketError, acceptSocket);
                    }, null);
                    
                    break;
                default:
                    // 异常错误                    
                    Env.L.Error($"socket error: {e.LastOperation}");
                    break;
            }
        }

        private void acceptAsync()
        {
            while (true)
            {
                this.innArgs.AcceptSocket = null;
                if (this._socket == null)
                    break;
                if (this._socket.AcceptAsync(this.innArgs))
                {
                    // 等待异步
                    return;
                }
                // 立刻返回
                OnAcceptComplete(this.innArgs.SocketError, this.innArgs.AcceptSocket);
            }
        }

        private void OnAcceptComplete(SocketError socketError, Socket acceptSocket)
        {
            Env.L.Info($"got net socket: {acceptSocket} thread: {Thread.CurrentThread.ManagedThreadId}");

            // 添加连接
            if(socketError != SocketError.Success)
            {
                if (acceptSocket != null)
                    SocketUtil.SafeClose(acceptSocket);
                Env.L.Error("error accept:" + socketError);
                acceptAsync();
                return;
            }

            _onAccept?.Invoke(acceptSocket);
            // next
            acceptAsync();
        }
    }
}

