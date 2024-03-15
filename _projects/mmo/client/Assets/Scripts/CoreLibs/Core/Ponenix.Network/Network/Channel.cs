using System;
using System.IO;
using Phoenix.Utils;

namespace Phoenix.Network
{
    // 提供一层收发缓冲
    public class Channel : IChannel, IDisposable
    {
        private int MAX_SEND_BUF_SIZE = 5 * 1024 * 1024;
        // 5ms
        private const int SOCKET_SEND_INTERVAL = 5;

        private CircularBuffer _recvBuf = new CircularBuffer();
        private CircularBuffer _sendBuf = new CircularBuffer();
        
        private ISession _session;
        private IProtocol _protocol;

        // 避免过度调用socket.send，控制频率        
        private LongIntervalCtrl _sendInterval = new LongIntervalCtrl();
        
        // TODO: 可以考虑IMsg的序列化也在pull时候调用

        public Channel(ISession session,
            IProtocol protocol)
        {
            _session = session;
            _protocol = protocol;
        }

        public int GetHandle()
        {
            return _session.GetHandle();
        }

        public void Clear()
        {
            _session = null;
            _protocol = null;
            Dispose();
        }

        public void StopSession()
        {
            if (_session == null)
                return;
            _session.ForceStop();
        }

        public void Dispose()
        {   
            _recvBuf.Dispose();
            _sendBuf.Dispose();
        }

        public void OnGotData(byte[] data, int offset, int count)
        {
            // 组织成消息
            // 接收由session串行控制，无需lock
            {
                _recvBuf.Write(data, offset, count);
                _protocol.MakeMsg(_recvBuf);
            }   
        }        

        public int PullSendData(byte[] buf, int offset, int max)
        {
            if (0 == _sendBuf.Length)
                return 0;
            lock (_sendBuf)
            {
                return _sendBuf.Read(buf, offset, max);
            }
        }

        public int GetSendDataSize()
        {
            return (int)_sendBuf.Length;
        }        

        public void SendMsg(Stream s)
        {
            if (s == null)
                return;
            lock (_sendBuf)
            {
                _sendBuf.Write(s);
            }

            trySendDataImmediately();
            tryStopSessionBufTooLarge();
        }

        public void SendMsg(byte[] bytes)
        {
            if (bytes == null)
                return;
            if (_session == null)
                return;
            lock (_sendBuf)
            {
                _sendBuf.Write(bytes, 0, bytes.Length);
            }

            trySendDataImmediately();
            tryStopSessionBufTooLarge();
        }

        private bool checkBufTooLarge()
        {
            if (_sendBuf.Length >= MAX_SEND_BUF_SIZE)
                return true;
            return false;
        }

        private bool tryStopSessionBufTooLarge()
        {
            if(checkBufTooLarge())
            {
                _session.Stop();
                return true;
            }
            return false;
        }

        public void Update()
        {
            trySendDataImmediately();
        }        
        
        private void trySendDataImmediately()
        {
            int size = GetSendDataSize();
            if (0 == size)
                return;
            // 数据不够大，时间间隔太低，则不马上发送
            if (size < 500 && !_sendInterval.CanDo(TimeUtil.HiNowMs(), SOCKET_SEND_INTERVAL))
                return;
            // call session
            _session.SendImmediately();
        }
    }

}

