
using System;
using System.Collections.Generic;
using System.Threading;

namespace Phoenix.Scheduler
{
    // 方便的获取Context
    // 并且同线程的异步await会保持在此线程唤醒执行下一步
    // 提供2个额外的工作线程

    // 主线程使用pull模式，所以要保证主线程的loop才能正常
    // await
    public class ThreadMgr
    {
        private static ThreadMgr _it;
        public static ThreadMgr it { get { return _it; } }
        private WorkThread _main = new WorkThread();

        List<WorkThread> _threads = new List<WorkThread>();
        private Random _rand = new Random();

        private ThreadMgrConfig _config = new ThreadMgrConfig();       

        public static void Start(int poolNum, ThreadMgrConfig config = null)
        {            
            _it = new ThreadMgr();
            
            if (config != null)
                _it._config = config;
            _it.start(poolNum);
        }

        // call in mainThread
        private void start(int poolNum)
        {            
            _main.StartWithPullMode(_config.skipNewMainContext);
            _threads.Add(_main);           

            // 启动多个线程
            for(var i = 0; i < poolNum; i ++)
            {
                var one = new WorkThread();
                one.StartAsWorkThread();

                // 等待启动完毕
                while(!one.IsStarted())
                {
                    Thread.Sleep(1);
                }
                _threads.Add(one);
            }
        }

        public static void Stop()
        {
            if (_it == null)
                return;
            _it.stop();
        }

        // 主线程的任务拉取执行
        public static void RunStep()
        {
            it.mainThread.PullUpdate();
        }

        private void stop()
        {
            for (var i = 0; i < _threads.Count; i++)
            {
                var one = _threads[i];
                one.Stop();
            }
        }

        public WorkThread mainThread { get { return _main; } }

        public SynchronizationContext mainContext { get { return _main.context; } }

        public SynchronizationContext GetContext(int threadId)
        {
            //lock(_threads)
            {
                for (var i = 0; i < _threads.Count; i++)
                {
                    var one = _threads[i];
                    if (one.threadId == threadId)
                        return one.context;
                }
            }
            return null;
        }

        // 获取当前Context
        public SynchronizationContext GetCurContext() 
        {
            return GetContext(Thread.CurrentThread.ManagedThreadId);
        }

        public SynchronizationContext GetOtherContext(int threadId)
        {   
            //lock (_threads)
            {
                // 随机下
                List<WorkThread> threads = new List<WorkThread>();
                for (var i = 0; i < _threads.Count; i++)
                {
                    var one = _threads[i];
                    if (one.threadId != threadId)
                        threads.Add(one);
                }

                if (threads.Count > 0)
                {
                    return threads[_rand.Next(0, threads.Count)].context;
                }
            }
            return null;
        }

        public SynchronizationContext GetOtherContext()
        {
            return GetOtherContext(Thread.CurrentThread.ManagedThreadId);
        }

        public SynchronizationContext GetWorkContext()
        {
            return GetOtherContext(_main.threadId);
        }

        // [0,)
        public SynchronizationContext GetWorkContextByIndex(int indexWork)
        {
            // 获取第几个
            var index = indexWork + 1;
            if (_threads.Count == 1)
            {
                // 归到主线程
                index = 0;
            }
            else
            {                
                if (index < 1 || index >= _threads.Count)
                    index = 1;
            }

            return _threads[index].context;
        }

        public WorkThread GetThreadByIndex(int indexWork)
        {
            // 获取第几个
            var index = indexWork + 1;
            if (_threads.Count == 1)
            {
                // 归到主线程
                index = 0;
            }
            else
            {
                if (index < 1 || index >= _threads.Count)
                    index = 1;
            }

            return _threads[index];
        }
    }
}
