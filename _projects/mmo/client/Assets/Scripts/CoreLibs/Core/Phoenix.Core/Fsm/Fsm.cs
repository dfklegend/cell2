
namespace Phoenix.Fsm
{
    /*
     * 状态机
     */
    public interface IFSMState
    {
        void SetCtrl(IFSMCtrl ctrl);
        void OnEnter();
        void OnLeave();
        void Update();
        bool IsOver();
    }

    public interface IFSMCtrl { }

    public delegate void UpdateDelegate(IFSMCtrl ctrl);

    public class BaseFSMCtrl<T> : IFSMCtrl
    {
        public IFSMFactory<T> factory;

        // 状态自动迁移
        public UpdateDelegate UpdateHandler;

        private T _stateType;
        private IFSMState _state;

        public T GetStateType() { return _stateType; }
        public IFSMState GetState() { return _state; }
        
        public void Update()
        {   
            _state?.Update();
            UpdateHandler?.Invoke(this);
        }

        public void ChangeState(T t)
        {
            if (_state != null && t.Equals(_stateType) )
                return;
            var newState = factory.CreateState(t);
            if (newState == null)
                return;
            newState.SetCtrl(this);
            _state?.OnLeave();

            _stateType = t;
            _state = newState;
            _state.OnEnter();
        }
    }

    public interface IFSMFactory<T>
    {
        IFSMState CreateState(T t);
        T GetInitState();
    }
}

