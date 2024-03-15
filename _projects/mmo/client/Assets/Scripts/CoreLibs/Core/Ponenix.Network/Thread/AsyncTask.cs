using System;
using System.Runtime.CompilerServices;
using System.Threading;

namespace Phoenix.Scheduler
{
    public delegate void AsyncTaskCallback(Action callback, object state);

    public class AsyncTask
    {  
        public bool IsFinished => false;

        // 执行完了之后, cb回来
        // 这样外部，可以包装成异步函数
        AsyncTaskCallback _task;       

        public AsyncTask(AsyncTaskCallback task)
        {
            _task = task;            
        }

        public void Do(object state)
        {
            if(_task != null)
            {
                _task.Invoke(() => { _awaiter.Next(); }, state);
            }
            
            Console.WriteLine($"Finish Thread:{Thread.CurrentThread.ManagedThreadId}");
            
        }

        AsyncTaskAwaiter _awaiter;
        public AsyncTaskAwaiter GetAwaiter()
        {            
            _awaiter = new AsyncTaskAwaiter(this);
            return _awaiter;
        }
    }

    public class AsyncTaskAwaiter : ICriticalNotifyCompletion
    {
        private readonly AsyncTask awaitable;
        private readonly SynchronizationContext capturedContext = SynchronizationContext.Current;
        private Action _continuation;

        public AsyncTaskAwaiter(AsyncTask awaitable) => this.awaitable = awaitable;
        public bool IsCompleted
        {
            get
            {
                return false;
            }
        }

        public void GetResult()
        {   
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
