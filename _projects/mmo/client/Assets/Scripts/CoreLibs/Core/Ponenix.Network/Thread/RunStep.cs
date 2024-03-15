
using System;

namespace Phoenix.Scheduler
{
    public interface IWorkThreadRunStep
    {
        void Update();
    }
    
    public class ActionRunStep: IWorkThreadRunStep
    {
        Action _action;
        public ActionRunStep(Action action)
        {
            _action = action;
        }

        public void Update()
        {
            _action();
        }
    }

    public class RunStep
    {
        public int id;
        public IWorkThreadRunStep step;
        // -1代表循环任务
        // 否则执行固定次数
        public int runTimes = -1;
    }

    // 方便使用
    public class RunStepCtrl
    {
        private WorkThread _thread;
        private int _stepId = -1;
        
        public RunStepCtrl(WorkThread thread)
        {
            _thread = thread;
        }

        public void Start(Action action, int runTimes = -1)
        {
            if (_stepId != -1)
                Stop();
            _stepId = _thread.AddRunStep(action, runTimes);
        }

        public void Stop()
        {
            if (_stepId == -1)
                return;
            _thread.RemoveRunStep(_stepId);
            _stepId = -1;
        }
    }
}
