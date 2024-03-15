using System;
using System.Runtime.CompilerServices;
using System.Threading;

namespace Phoenix.Network.Protocol.Pomelo
{
    // Request async支持
    public class RequestTask : ICriticalNotifyCompletion
    {
        private readonly SynchronizationContext capturedContext = SynchronizationContext.Current;
        private Action _continuation;
        RequestResult _result;

        public RequestTask(PomeloNode node, string route, byte[] data)
        {
            node.Request(route, data, (result) => {
                _result = result;
                Next(); 
            });
        }

        public RequestTask GetAwaiter()
        {
            return this;
        }

        public bool IsCompleted
        {
            get
            {
                return false;
            }
        }

        public RequestResult GetResult()
        {
            return _result;
        }

        public void OnCompleted(Action continuation)
        {
            _continuation = continuation;
        }

        public void UnsafeOnCompleted(Action continuation)
        {
            _continuation = continuation;
        }

        public void Next()
        {
            if (capturedContext != null)
            {
                var continuation = _continuation;
                capturedContext.Post(state => continuation(), null);
            }
            else
                _continuation();
        }

    }
}

