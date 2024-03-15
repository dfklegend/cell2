
using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;

namespace Phoenix.Scheduler
{
    public class OneOper
    {
        public SendOrPostCallback cb;
        public object state;
        public ManualResetEvent wait;
   
        public AsyncTask asyncTask;

        public OneOper(SendOrPostCallback cb, object state)
        {
            this.cb = cb;
            this.state = state;
        }

        public OneOper(SendOrPostCallback cb, object state, ManualResetEvent wait)
        {
            this.cb = cb;
            this.state = state;
            this.wait = wait;
        }

        public static OneOper NewAsync(AsyncTask task, object state)
        {
            var obj = new OneOper(null, state);
            obj.asyncTask = task;        
            return obj;
        }
    }

    // 同步上下文
    // . 可以主动用Post来将任务推送到对应线程
    // . 确保await在调用线程内触发(会自动使用当前的Context)
    //   在线程内SynchronizationContext.SetSynchronizationContext( ... );
    public class ThreadSynchronizationContext: SynchronizationContext
    {
        private int _threadId = -1;
        Queue<OneOper> _opers = new Queue<OneOper>();
        List<Task> _tasks = new List<Task>();

        public ThreadSynchronizationContext(int threadId = -1)
        {
            _threadId = threadId;
        }

        public void BindThread(int threadId)
        {
            _threadId = threadId;
        }

        private bool isThreadBind()
        {
            return _threadId != -1;
        }

        public override void Post(SendOrPostCallback cb, object state)
        {
            if (_threadId == Thread.CurrentThread.ManagedThreadId)
            {
                // 已经是此线程                
                //Console.WriteLine("same");
            }

            lock (_opers)
            {
                _opers.Enqueue(new OneOper(cb, state, null));
            }
        }

        public void PostAlways(SendOrPostCallback cb, object state)
        {
            lock (_opers)
            {
                _opers.Enqueue(new OneOper(cb, state, null));
            }
        }

        // 调用会等待执行完毕
        public override void Send(SendOrPostCallback cb, object state)
        {
            if (_threadId == Thread.CurrentThread.ManagedThreadId)
            {
                Network.Env.L.Warning("Warnning, SynchronizationContext.Send in same thread!");
                // 已经是此线程，直接执行
                cb(state);
                return;
            }

            // 等待执行完毕
            var wait = new ManualResetEvent(false);
            lock (_opers)
            {
                _opers.Enqueue(new OneOper(cb, state, wait));
            }
            wait.WaitOne();
            wait.Dispose();
        }

        /*
         * 使用
         * 让工作线程完成之后，再交回
         * await PostAsync((cb)=> {
         *  await something
         *  ...
         *  cb()
         * }, some)
         */
        public AsyncTask PostAsync(AsyncTaskCallback cb, object state)
        {
            var task = new AsyncTask(cb);
            lock (_opers)
            {
                _opers.Enqueue(OneOper.NewAsync(task, state) );
            }
            return task;
        }

        // 由触发线程调用
        public void ExecuteTasks()
        {
            if (!isThreadBind())
                return;
            lock (_opers)
            {
                while (true)
                {
                    var head = tryPop();
                    if (head == null)
                        break;
                    doOper(head);
                }
            }
        }

        private OneOper tryPop()
        {
            if (_opers.Count == 0)
                return null;
            return _opers.Dequeue();
        }

        private void doOper(OneOper oper)
        {
            try
            {
                implDoOper(oper);
            }
            catch(Exception e)
            {
                Network.Env.L.Error("ThreadSynchronizationContext.doOper Exception:");
                Network.Env.L.Error(e.ToString());
                Utils.SystemUtil.LogHandledException(e);
            }
        }

        private void implDoOper(OneOper oper)
        {
            if (oper.asyncTask != null)
            {
                oper.asyncTask.Do(oper.state);
                return;
            }
            oper.cb(oper.state);
            if (oper.wait != null)
                oper.wait.Set();
        }

        private void checkTasks()
        {
            if (_tasks.Count == 0)
                return;            
                    
        }
    }
}
