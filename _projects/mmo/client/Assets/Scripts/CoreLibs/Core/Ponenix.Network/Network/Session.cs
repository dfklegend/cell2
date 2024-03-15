using System;
using System.Net.Sockets;
using System.Threading;

namespace Phoenix.Network
{
    // 只负责读写
    // . 接收数据并Push到上层
    // . 从上层拉取数据发送
    public class TCPSession : ISession
    {
        private const int MAX_BUFFER_SIZE = 1500;
		private const int HALF_MAX_BUFFER_SIZE = MAX_BUFFER_SIZE/2;

		private long _id;
		public long id { get { return _id; } }
        Socket _socket;
        
        private SocketAsyncEventArgs _readArgs = new SocketAsyncEventArgs();
        private SocketAsyncEventArgs _writeArgs = new SocketAsyncEventArgs();

		// 保证只能启动一次ReceiveAsync
		//int _reading = 0;
		private byte[] _readBuf = new byte[MAX_BUFFER_SIZE]; 

		// 发送buf，发完再获取
		private byte[] _writeBuf = new byte[MAX_BUFFER_SIZE];
		// 后备buf，
		private byte[] _backBuf = new byte[MAX_BUFFER_SIZE];
		private int _writeOffset = 0;
		private int _writeSize = 0;
		
		// 只有一个SendAsync调用
		private int _writeToken = 0;

		private Action<long, SocketError> _cbOnError;

		private bool _started = false;
		private bool _error = false;

		public NetConfig config;
		IChannel _channel;	

        public TCPSession(long id, Socket socket, Action<long, SocketError> cb)
        {
			_id = id;
            _socket = socket;
			_cbOnError = cb;

			this._readArgs.Completed += this.OnRecvComplete;						
			this._readArgs.SetBuffer(_readBuf, 0, MAX_BUFFER_SIZE);

			this._writeArgs.Completed += this.OnSendComplete;
		}

		public int GetHandle()
        {
			if (_socket == null)
				return 0;
			return _socket.Handle.ToInt32();

		}

		public void SetCBError(Action<long, SocketError> cb)
        {
			_cbOnError = cb;
        }

		public void SetChannel(IChannel channel)
        {
			_channel = channel;
        }

		// 开始收发数据
		public void Start()
        {
			if (_started)
				return;			

			_started = true;
			StartRecv();
			StartSend();
        }

		public void Stop()
        {
			if (_socket == null)
				return;
			Console.WriteLine($"{id} stop");
			SocketUtil.SafeClose(_socket);
			_socket = null;
			_cbOnError = null;
			_channel.Clear();
			_channel = null;
			_readArgs.Dispose();
			_writeArgs.Dispose();
        }

		public void ForceStop()
        {
			var cbError = _cbOnError;
			Stop();
			cbError?.Invoke(_id, SocketError.TimedOut);
        }

		private void StartRecv()
		{			
			reqRecv();
		}

		private void reqRecv()
		{
			while (true)
			{
				if (this._socket == null)
					return;
                try
                {
					if (this._socket.ReceiveAsync(this._readArgs))
					{
						return;
					}
				}
				catch(Exception e)
                {
					//if(this._socket != null)
					{
						Console.WriteLine($"{id} Exception:");
						Console.WriteLine(e);
						Utils.SystemUtil.LogHandledException(e);
					}
					
					return;
                }
				
				if (!this.HandleRecv(this._readArgs))
					break;
			}
		}

		private void OnRecvComplete(object sender, SocketAsyncEventArgs e)
		{
			if (!this.HandleRecv(e))
				return;

			if (this._socket == null)
			{
				return;
			}
			this.StartRecv();
		}

		// return false if error
		private bool HandleRecv(SocketAsyncEventArgs e)
		{
			if (this._socket == null)
			{
				return false;
			}
			
			if (e.SocketError != SocketError.Success)
			{
				onError(e.SocketError);
				return false;
			}

			if (e.BytesTransferred == 0)
			{
				//Console.WriteLine($"{_socket.Handle.ToInt32()} recv 0 bytes");		
				onError(SocketError.ConnectionAborted);
				return false;
			}

			//Console.WriteLine($"{_socket.Handle.ToInt32()} recv {e.BytesTransferred} bytes");
			//Env.L.FileLog($"{_socket.Handle.ToInt32()} recv {e.BytesTransferred} bytes");
			//Env.L.FileLogBuf($"{GetHandle()} {e.BytesTransferred} ", _readBuf, 0, e.BytesTransferred, 10);

			//Console.WriteLine($"{_id} recved: {e.BytesTransferred} Thread:{Thread.CurrentThread.ManagedThreadId}");
			_channel.OnGotData(_readBuf, 0, e.BytesTransferred);
			return true;
		}

		private void onError(SocketError err)
        {
			_error = true;
			Env.L.Info($"{_id} error: {err}");
			_cbOnError?.Invoke(_id, err);
		}

		private bool tryReqWriteToken()
        {
			if (1 == Interlocked.Increment(ref _writeToken))
				return true;
			// 还回去
			Interlocked.Decrement(ref _writeToken);
			return false;
		}

		private void backWriteToken()
		{	
			// 还回去
			Interlocked.Decrement(ref _writeToken);			
		}

		public void StartSend()
        {
			if (!tryReqWriteToken())
				return;
			doSend();
        }
	
		public void SendImmediately()
        {
			if (_error)
				return;
			
			if (!tryReqWriteToken())
				return;
			// 去线程里面调用
			ThreadPool.QueueUserWorkItem((state) => {
				doSend();
			});

			//doSend();
		}

		private bool hasPendingSendData()
        {
			return _channel.GetSendDataSize() > 0;			
        }

		private void doSend()
		{				 
			while (true)
			{
				try
				{
					if (this._socket == null)
					{
						return;
					}					

					tryMakeWriteBuf();
					if(_writeSize == 0)
                    {
						// 并没有数据需要发送						
						// 等待外部send触发
						backWriteToken();

						// 由新的send触发，极端如下情况
						/*
						 * SendThread.tryMakeWriteBuf
						 * _writeSize == 0
						 * 
						 * SwitchThread -> CallThread
						 * CallThread.SendImmediately 
						 *		tryReqWriteToken false
						 *		
						 * SwitchThread -> SendThread
						 * SendThread.backWriteToken						 
						 * 
						 */
						if (hasPendingSendData())
						{
							if (tryReqWriteToken())
								// 说明刚好错开了
								// do send again
								continue;
							else
								return;
						}
						else
                        {
							return;
                        }
					}

					if (this._error)
						return;
                    try
                    {
						_writeArgs.SetBuffer(_writeBuf, _writeOffset, _writeSize);
						//Env.L.FileLog($"{GetHandle()} begin send {_writeSize}");
						//Env.L.FileLogBuf($"{GetHandle()}", _writeBuf, _writeOffset, _writeSize, 10);
						if (this._socket.SendAsync(this._writeArgs))
						{
							return;
						}
					}
					catch(Exception e)
                    {
						//if(this._socket != null)
						{
							Console.WriteLine($"{id} Exception:");
							Console.WriteLine(e);
							Utils.SystemUtil.LogHandledException(e);
						}
						return;
					}					

					if (!HandleSend(this._writeArgs))
						break;
				}
				catch (Exception e)
				{	
					Console.WriteLine($"doSend Exception:");
					Console.WriteLine(e);
					Utils.SystemUtil.LogHandledException(e);
				}
			}
		}

		private void tryMakeWriteBuf()
        {
			if (!hasPendingSendData())
				return;
			// 重新整理buf
			if(_writeOffset > HALF_MAX_BUFFER_SIZE / 2)
            {
				//Env.L.FileLog($"{GetHandle()} swapBuf");
				//Env.L.FileLogBuf($"{GetHandle()} before", _writeBuf, _writeOffset, _writeSize, 10);
				swapBuf();
				//Env.L.FileLogBuf($"{GetHandle()} after", _writeBuf, _writeOffset, _writeSize, 10);
			}
			int bufRestSize = MAX_BUFFER_SIZE - _writeSize - _writeOffset;			
			int got = _channel.PullSendData(_writeBuf, _writeOffset, bufRestSize);
			_writeSize += got;
		}

		private void swapBuf()
        {
			Array.Copy(_writeBuf, _writeOffset, _backBuf, 0, _writeSize);
			
			byte[] temp;
			temp = _backBuf;
			_backBuf = _writeBuf;
			_writeBuf = temp;

			_writeOffset = 0;
        }

		private void OnSendComplete(object sender, SocketAsyncEventArgs e)
		{
			if (!HandleSend(e))
				return;
			this.doSend();
		}

		private bool HandleSend(SocketAsyncEventArgs e)
		{		
			if (e.SocketError != SocketError.Success)
			{
				this.onError(e.SocketError);
				return false;
			}

			if (e.BytesTransferred == 0)
			{
				Env.L.Warning($"{_socket.Handle.ToInt32()} send 0 bytes");					
			}

			//Env.L.FileLog($"{_socket.Handle.ToInt32()} send {e.BytesTransferred} bytes");
			//Console.WriteLine($"{_socket.Handle.ToInt32()} send {e.BytesTransferred} bytes");

			_writeOffset += e.BytesTransferred;
			_writeSize -= e.BytesTransferred;
			if (_writeSize == 0)
				_writeOffset = 0;
			return true;
		}
	}
}

