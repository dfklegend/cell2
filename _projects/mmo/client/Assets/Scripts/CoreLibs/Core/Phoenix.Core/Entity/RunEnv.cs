using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Phoenix.Core
{   
    /*
     * 可以获取root节点 
     * AddComponent来增加新的模块
     * 
     * 一个RunEnv只运行在一个线程内
     * 考虑一个环境提供一个线程安全的公共模块
     */
    public class RunEnv
    {
        private Entity _root;
        public Entity root { get { return _root; } }

        // 一些和运行环境相关的系统模块功能
        private TimerManager _timer;
        public TimerManager timer { get { return _timer; } }

        private EventCenter<int> _events = new EventCenter<int>();
        public EventCenter<int> events { get { return _events; } }

        private ComponentPoolMgr _pools = new ComponentPoolMgr();
        public ComponentPoolMgr pools { get { return _pools; } }

        private ComponentRunMgr _componentRun = new ComponentRunMgr();
        public ComponentRunMgr componentRun { get { return _componentRun; } }

        private EntityMgr _entites = new EntityMgr();
        public EntityMgr entites { get { return _entites; } }

        private SceneMgr _scenes;
        public SceneMgr scenes { get { return _scenes; } }

        public void CreateEnv()
        {
            _root = new Entity();
            _root.SetId(AppEnv.AllocUniqueId());
            _root.SetEnv(this);
            _entites.AddEntity(_root);

            var timer = _root.AddComponent<TimerMgrComponent>();
            _timer = timer.timer;

            _scenes = SceneMgr.Create(this);
        }

        public void Update()
        {
            updateSystemModules();
            
            _root.Update();

            // batch模式执行所有的components
            _componentRun.Update();
            // 判断回收
            _componentRun.RecycleDestroyed(_pools);
        }

        private void updateSystemModules()
        {   
        }

        public Entity CreateEntity(Entity parent)
        {
            if (parent == null)
                parent = _root;
            var e = new Entity();
            e.SetId(AppEnv.AllocUniqueId());
            e.parent = parent;
            e.SetEnv(parent.env);
            _entites.AddEntity(e);
            return e;
        }

        public void RemoveEntity(Entity e)
        {
            _entites.RemoveEntity(e);
        }

        public Entity GetEntity(ulong id)
        {
            return _entites.GetEntity(id);
        }

        public Entity FindEntity(string path)
        {
            return _root.Find(path);
        }
    }    

    public class TimerMgrComponent : BaseComponent
    {
        // 一些和运行环境相关的系统模块功能
        private TimerManager _timer = new TimerManager();
        public TimerManager timer { get { return _timer; } }

        public override int GetExecuteOrder()
        {
            return ComponentOrderDefine.SYSTEM;
        }

        public override void OnUpdate()
        {
            base.OnUpdate();
            _timer.Expire();
        }
    }
}
