
using System;
using System.Collections.Generic;
using System.Threading;

namespace Phoenix.Scheduler
{
    // 用户线程，拥有context，可以使用Post推送到对应线程执行
    public class WorkThread
    {
        private ThreadSynchronizationContext _context = new ThreadSynchronizationContext();
        private SynchronizationContext _externContext;

        private Thread _thread;
        private int _threadId = -1;
        public int threadId { get { return _threadId; } }        
        
        private List<RunStep> _steps = new List<RunStep>();
        private int _nextStepId = 0;


        private volatile bool _stopped = false;
        public SynchronizationContext context 
        { 
            get 
            { 
                
                return _externContext!= null?_externContext:_context; 
            } 
        }

        public WorkThread() {}
       
        public void StartAsWorkThread()
        {
            _thread = new Thread(threadFunc);
            _thread.Start();
        }

        // 使用当前线程拉取(一般主线程)
        public void StartWithPullMode(bool skipNewMainContext)
        {
            bindCurThreadContext(skipNewMainContext);
        }
        

        public bool IsStarted()
        {
            return _threadId != -1;
        }

        public void Stop()
        {
            _stopped = true;            
        }

        private void threadFunc()
        {
            bindCurThreadContext(false);
            loopUpdate();
        }

        private void bindCurThreadContext(bool skipNewMainContext)
        {
            _threadId = Thread.CurrentThread.ManagedThreadId;
            Network.Env.L.Info($"bind context {_threadId}");
            // for unity
            if (skipNewMainContext)
            {
                // must not be null                
                if (SynchronizationContext.Current == null)
                {
                    Network.Env.L.Info($"cur context is null");
                    throw new Exception("cur context is null when use skipNewMainContext");
                }
                _externContext = SynchronizationContext.Current;
                return;
            }          

            _context.BindThread(_threadId);
            SynchronizationContext.SetSynchronizationContext(_context);
        }

        private void loopUpdate()
        {
            while(!_stopped)
            {
                update();
                // 能否用唤醒机制
                Thread.Sleep(1);
            }
        }

        private void update()
        {
            updateSteps();
            _context.ExecuteTasks();
        }

        private void updateSteps()
        {
            if (_steps.Count == 0)
                return;
            lock (_steps)
            {
                for (var i = 0; i < _steps.Count;)
                {
                    var one = _steps[i];

                    if(one.runTimes > 0 || one.runTimes == -1)
                    {
                        one.step.Update();

                        if (one.runTimes > 0)
                        {
                            one.runTimes--;
                            if (one.runTimes == 0)
                            {
                                // skip add i
                                _steps.RemoveAt(i);
                                continue;
                            }
                        }                        
                    }

                    i++;
                }                
            }
        }

        public void PullUpdate()
        {
            update();
        }


        // runStep
        private int allocRunStepId()
        {
            return Interlocked.Increment(ref _nextStepId);
        }

        // 插入来增加处理
        public int AddRunStep(IWorkThreadRunStep step, int runTimes = -1)
        {
            if (runTimes == 0)
                return -1;
            lock (_steps)
            {
                var one = new RunStep();
                one.id = allocRunStepId();
                one.step = step;
                one.runTimes = runTimes;

                _steps.Add(one);
                return one.id;
            }
        }

        public int AddRunStep(Action action, int runTimes = -1)
        {
            return AddRunStep(new ActionRunStep(action), runTimes);
        }

        public int findIndex(Func<RunStep, bool> func)
        {
            for (var i = 0; i < _steps.Count; i++)
            {
                var one = _steps[i];
                if (func(one))
                    return i;
            }
            return -1;
        }

        public void RemoveRunStep(IWorkThreadRunStep step)
        {
            lock (_steps)
            {
                var index = findIndex((one) => { return one.step == step; });
                if (index == -1)
                    return;
                _steps.RemoveAt(index);
            }
        }

        public void RemoveRunStep(int id)
        {
            lock (_steps)
            {
                var index = findIndex((one) => { return one.id == id; });
                if (index == -1)
                    return;
                _steps.RemoveAt(index);
            }
        }

    }
}
