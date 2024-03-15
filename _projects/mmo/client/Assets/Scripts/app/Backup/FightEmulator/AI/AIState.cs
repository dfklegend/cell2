using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;
using Phoenix.Fsm;
using Phoenix.Utils;

namespace Phoenix.Game.FightEmulator
{  
    public abstract class BaseAIState : IFSMState
    {
        protected AIStateCtrl _ctrl;
        public void SetCtrl(IFSMCtrl ctrl) 
        {
            _ctrl = ctrl as AIStateCtrl;
        }
        public abstract void OnEnter();
        public abstract void OnLeave();
        public abstract void Update();
        public virtual bool IsOver() { return false; }
    }

    public class AIStateFactory : Singleton<AIStateFactory>, IFSMFactory<int>
    {
        IntToClassFactory<BaseAIState> _factory = new IntToClassFactory<BaseAIState>();

        public AIStateFactory()
        {
            _factory.RegisterAll();
        }

        public IFSMState CreateState(int t)
        {
            return _factory.Create(t);
        }

        public int GetInitState()
        {
            return (int)eAIState.Init;
        }
    }

    public class AIStateCtrl : BaseFSMCtrl<int>
    {
        private Character _owner;
        public Character owner { get { return _owner; } }
        public void SetCharacter(Character c) { _owner = c; }

        private AICtrl _ai;
        public AICtrl ai { get { return _ai; } }
        public void SetAI(AICtrl ai) { _ai = ai; }

        public static AIStateCtrl New()
        {
            var ctrl = new AIStateCtrl();
            ctrl.factory = AIStateFactory.It;
            return ctrl;
        }
    }
}// namespace Phoenix
