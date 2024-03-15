using Phoenix.Fsm;
using Phoenix.Utils;
using Phoenix.Core;

namespace Phoenix.Game
{   
    public abstract class BaseAppState : IFSMState
    {
        public void SetCtrl(IFSMCtrl ctrl) { }
        public abstract void OnEnter();
        public abstract void OnLeave();
        public abstract void Update();
        public virtual bool IsOver() { return false; }
    }

    public class AppStateFactory : Singleton<AppStateFactory>, IFSMFactory<int>
    {
        IntToClassFactory<BaseAppState> _factory = new IntToClassFactory<BaseAppState>();

        public AppStateFactory()
        {
            _factory.RegisterAll();
        }

        public IFSMState CreateState(int t)
        {
            return _factory.Create(t);
        }

        public int GetInitState()
        {
            return (int)eAppState.Init;
        }
    }

    public class AppStateCtrl : BaseFSMCtrl<int>
    {
    }

    public static class FSMFactory
    {
        public static AppStateCtrl Create()
        {
            var ctrl = new AppStateCtrl();
            ctrl.factory = AppStateFactory.It;
            return ctrl;
        }
    }
}

