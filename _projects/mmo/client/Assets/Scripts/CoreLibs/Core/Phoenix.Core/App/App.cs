using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Phoenix.Core
{
    public enum eAppState
    {
        Init = 0,
        WaitReady,
        Ready,
        Stopping,
        Stopped
    }

    public abstract class BaseApp
    {
        private eAppState _state = eAppState.Init;
        
        private void setState(eAppState state)
        {
            _state = state;
        }

        public void Prepare()
        {
            OnPrepare();
        }

        // 注册AppComponent
        public abstract void OnPrepare();

        public void Start()
        {
            AppComponentMgr.StartAll();
            OnStart();
            _state = eAppState.WaitReady;
        }

        public abstract void OnStart();


        // 所有component ready了之后
        public virtual void OnReady()
        {

        }

        public bool IsReady()
        {
            return _state == eAppState.Ready;
        }

        // 需要外部调用
        public void Update()
        {
            AppComponentMgr.Update();
            updateStates();
            OnUpdate();
        }

        private void updateStates()
        {
            if (_state == eAppState.WaitReady)
            {
                if (AppComponentMgr.IsAllReady())
                {
                    setState(eAppState.Ready);
                    OnReady();
                }
                return;
            }            
            if(_state == eAppState.Stopping)
            {
                if(AppComponentMgr.IsAllStopped())
                {
                    OnStop();
                    setState(eAppState.Stopped);                    
                    PConsole.Log("App.Stopped");                    
                }
            }
        }

        public abstract void OnUpdate();

        public void Stop()
        {
            PConsole.Log("App.Stop");
            AppComponentMgr.StopAll();
            setState(eAppState.Stopping);
        }

        public abstract void OnStop();

        // 所有的stop完毕
        public bool CanSafeExit()
        {
            return eAppState.Stopped == _state;
        }
    }
}
