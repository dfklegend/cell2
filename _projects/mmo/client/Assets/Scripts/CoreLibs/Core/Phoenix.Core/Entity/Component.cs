using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Phoenix.Core
{   
    public class ComponentOrderDefine
    {
        // 0-1000，系统级
        public const int SYSTEM = 1000;
        // >=2000 用户定义
        public const int NORMAL = 2000;
        
    }

    public class BaseComponent
    {
        private Entity _owner;
        public Entity owner { get { return _owner; } }

        private bool _destroyed = false;
        public bool destroyed { get { return _destroyed; } }

        private int _updateTimes = 0;

        // 排序 决定控件运行顺序
        // 小的越优先执行
        public virtual int GetExecuteOrder() { return ComponentOrderDefine.NORMAL; }
        // 如果很明确的不需要update驱动
        public virtual bool ForceNoUpdate() { return false; }
        
        public virtual void Create(params object[] args) { }
        
        public void SetOwner(Entity owner)
        {
            _owner = owner;
        }


        public virtual void OnCreate() { }

        // 首次被调用
        public virtual void OnAwake() { }
        public virtual void OnUpdate() { }

        public void Update()
        {
            if (_updateTimes == 0)
                OnAwake();
            else
                OnUpdate();
            _updateTimes++;
        }

        // 做清理
        public virtual void OnDestroy() { }

        public void Destroy()
        {
            if (_destroyed)
                return;
            if (_owner == null)
                return;
            _owner.DestroyComponent(this);
        }

        // 只允许Entity调用
        public void DoDestroy() 
        {
            if (_destroyed)
                return;
            _destroyed = true;
            OnDestroy();
        }

        // 恢复成缺省状态，便于复用
        public virtual void Recycle()
        {
            _updateTimes = 0;
            _owner = null;
            _destroyed = false;
        }
    }
}
