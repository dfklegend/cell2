using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;

namespace Phoenix.Game
{	

    public interface ILogicImpl
    {
        void Start();
        void Update();
        void Destroy();
    }

    public class GameLogicCtrl : AppTask.BaseCommonTask
    {
        private static GameLogicCtrl _it;
        public static GameLogicCtrl It { get { return _it; } }
        private ILogicImpl _impl;

        static GameLogicCtrl()
        {
            _it = new GameLogicCtrl();
        }

        public void Init()
        {
            AppTask.CommonTaskDriver.It.AddTask(this);
        }

        public override void Update() 
        {
            _impl?.Update();
        }

        public void ChangeImpl<T>()
            where T: ILogicImpl
        {
            DestroyImpl();
            _impl = (T)Activator.CreateInstance(typeof(T));
            _impl.Start();
        }

        public void DestroyImpl()
        {
            if (_impl == null)
                return;
            _impl.Destroy();
            _impl = null;
        }

        public void Clear()
        {
            DestroyImpl();            
        }
    }
} // namespace Phoenix
