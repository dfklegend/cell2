using System;
using System.Net;
using System.Net.Sockets;
using System.Threading;
using Phoenix.Scheduler;

namespace Phoenix.Network
{
    // 发起连接
    public class TCPClient
    {
		long _id;
		public long id { get { return _id; } }
		private SynchronizationContext _main;
		NetConfig _config;

		// 连接成功之后，此socket会转交给Session对象
		private Socket _socket;
        private SocketAsyncEventArgs _connArgs = new SocketAsyncEventArgs();
		IPEndPoint _connectTar;

		// valid when connectted
		TCPConnection _conn;

		public Action<SocketError> cbConnectFail;
		public Action<SocketError> cbConnectionErr;
		public Action<TCPClient> cbConnectted;
		public Action<TCPClient> cbReady;
		public Action<TCPConnection, IMsg> cbProcessMsg;		


		public TCPClient(SynchronizationContext main, NetConfig config)
        {			
			_main = main;
			this._config = config;
			_connArgs.Completed += OnConnectComplete;			
		}

		// ip, port
		public void Connect(string ip, int port)
        {
			var ipEndPoint = new IPEndPoint(IPAddress.Parse(ip), port);
			Connect(ipEndPoint);
		}

		private bool hasPrevSocketNeedRelease()
        {
			return _socket != null || _conn != null;

		}

		public void Connect(IPEndPoint ipEndPoint)
        {
			if(hasPrevSocketNeedRelease())
            {
				doStop();
            }
            this._socket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            this._socket.NoDelay = true;
			_connectTar = ipEndPoint;

			_connArgs.RemoteEndPoint = ipEndPoint;			
			ConnectAsync();
			Env.L.Debug($"connect to {ipEndPoint}");
        }

		// 使用上次地址重新连接
		public void ReConnect()
        {			
			doStop();			
			Connect(_connectTar);
        }

		public void Stop()
        {		
			doStop();
        }

		public void doStop()
        {
			Env.L.Info($"TCPClient.doStop");
			if (_socket != null)
            {
				_socket.Close();
				_socket = null;
			}
			if (_conn != null)
            {
				_conn.Stop();
				_conn = null;
			}			
        }

		public bool IsConnected()
        {
			return _conn != null && _conn.IsReady();
        }

		public T GetProtocol<T>()
			where T : class, IProtocol
		{
			return _conn.GetProtocol<T>();
		}

		private void ConnectAsync()
		{
			try
			{
				if (this._socket.ConnectAsync(this._connArgs))
				{
					return;
				}
			}
			catch(Exception e)
            {
				Env.L.Error(e.ToString());
				Env.L.Error("Connect Exception:");
				Utils.SystemUtil.LogHandledException(e);
			}

			_main.Post((state) => { processConnectComplete(this._connArgs); }, null);			
		}

		private void OnConnectComplete(object sender, SocketAsyncEventArgs e)
		{
			_main.Post((state) => { processConnectComplete(this._connArgs); }, null);
		}

		private void processConnectComplete(SocketAsyncEventArgs e)
        {
			if (this._socket == null)
			{
				return;
			}
			if (e.SocketError != SocketError.Success)
			{
				Env.L.Error($"Connect error: {e.SocketError}");
				onConnectFail(e.SocketError);
				return;
			}

			// create 
			_id = ConnIDAllocer.Alloc();
			_conn = TCPBuilder.Build(_config, _id, this._socket);
			_conn.SetCBProcessMsg(cbProcessMsg);			
			
			_conn.session.SetCBError(onSessionErr);
			// 开始收消息
			_conn.Start();

			_conn.protocol.SetCBOnReady(() => {
				onReady();
			});

			onConnectted();			
			// 交给session处理了
			this._socket = null;
		}

		private void onConnectFail(SocketError error)
        {
			cbConnectFail?.Invoke(error);
        }

		private void onConnectted()
        {
			cbConnectted?.Invoke(this);
        }

		private void onReady()
        {
			_main.Post((state) => {
				cbReady?.Invoke(this);
			}, null);
		}

		// 由使用者驱动
		public void Update()
        {
			if (_conn == null)
				return;
			_conn.Update();
        }		

		private void onSessionErr(long id, SocketError err)
        {
			Env.L.Info($"TCPClient.onSessionErr {id} Thread:{Thread.CurrentThread.ManagedThreadId}");
			_main.Post((state) => {
				this.onErr(id, err);
			}, null);
		}

		private void onErr(long id, SocketError err)
        {
			Env.L.Info($"TCPClient.onErr {id} Thread:{Thread.CurrentThread.ManagedThreadId}");
			if (_conn == null)
				return;
			_conn.Stop();
			_conn = null;

			cbConnectionErr?.Invoke(err);
		}

		public IProtocol GetChannel()
		{
			if (_conn == null)
				return null;
			return _conn.protocol;
		}
	}
}

