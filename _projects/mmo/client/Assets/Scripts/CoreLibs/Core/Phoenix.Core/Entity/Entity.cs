using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Phoenix.Core
{
    /*
     * 
     */
    // Entity负责集成Component
    // Entity可以组织成节点树
    // Entity的Update有根节点驱动

    // 降低复杂度,
    // Entity Destroy
    //      Destroy立刻生效
    //      在下一次父节点Update从节点树中移除
    public partial class Entity
    {
        public delegate void OnEntityDestroy(Entity e);
        public static bool componentRunBatch = true;
        // true: 由ComponentRunMgr驱动component
        // false: 按节点顺序
        // TODO: 最后选择一种

        private ulong _id = 0;
        public ulong id { get { return _id; } }
        public void SetId(ulong id) { _id = id; }

        private List<BaseComponent> _comps = new List<BaseComponent>();

        private Entity _parent = null;
        private List<Entity> _childs = new List<Entity>();

        private bool _destroyed = false;
        public bool destroyed { get { return _destroyed; } }

        //public event OnEntityDestroy onDestroy;

        public T AddComponent<T>(params object[] args)
            where T: BaseComponent, new()
        {
            if (_destroyed)
                return null;
            var c = allocComponent<T>();
            // 初始化参数
            c.Create(args);
            
            c.SetOwner(this);
            _comps.Add(c);

            c.OnCreate();
            
            if(componentRunBatch)
            {
                env.componentRun.Add(c);
            }
            return c;
        }

        private T allocComponent<T>()
            where T : BaseComponent, new()
        {
            return _env.pools.Alloc<T>();
        }

        private void freeComponent(BaseComponent c)
        {
            if(componentRunBatch)
            {
                // will releas in runlist
                return;
            }

            c.Recycle();
            _env.pools.Free(c);
        }

        // call by Component.Destroy
        public void DestroyComponent(BaseComponent c)
        {
            if (c.destroyed)
                return;
            var index = _comps.IndexOf(c);
            if (index == -1)
                return;

            c.DoDestroy();           
        }        

        public T GetComponent<T>()
            where T: BaseComponent
        {
            var type = typeof(T);
            for(var i = 0; i < _comps.Count; i ++)
            {
                var one = _comps[i];
                if (one.destroyed)
                    continue;
                if (one is T)
                    return one as T;
            }
            return default;
        }

        public List<T> GetComponents<T>()
            where T : BaseComponent
        {
            List<T> result = new List<T>();
            var type = typeof(T);
            for (var i = 0; i < _comps.Count; i++)
            {
                var one = _comps[i];
                if (one.destroyed)
                    continue;
                if (one is T)
                {
                    result.Add(one as T);
                }                    
            }
            return result;
        }

        // TODO: 确认
        // 更新方式
        // 类似ECS组织
        // ComponentPool能否带来性能提升
        public void Update()
        {   
            if(!componentRunBatch)
            {
                // batch模式下
                // 由RunEnv统一驱动
                updateComponents();
                removeComponentsDestroyed();
            }
            
            updateChilds();
            removeEntityDestroyedOrNotChild();
        }

        private void updateComponents()
        {           
            for (var i = 0; i < _comps.Count; i++)
            {
                var one = _comps[i];
                if (one.destroyed)
                {                    
                    continue;
                }

                one.Update();
            }
        }

        private void removeComponentsDestroyed()
        {            
            for (var i = 0; i < _comps.Count;)
            {
                var one = _comps[i];
                if (one.destroyed)
                {
                    freeComponent(one);
                    _comps.RemoveAt(i);
                    continue;
                }

                i++;
            }
        }

        

        private bool IsValidChild(Entity child)
        {
            if (child.destroyed || child.parent != this)
                return false;
            return true;
        }

        private void updateChilds()
        {
            for (var i = 0; i < _childs.Count; i++)
            {
                var one = _childs[i];
                if (IsValidChild(one))
                {
                    one.Update();
                }
            }
        }
        
        private void removeEntityDestroyedOrNotChild()
        {
            for (var i = 0; i < _childs.Count;)
            {
                var one = _childs[i];
                if (!IsValidChild(one))
                {
                    _childs.RemoveAt(i);
                    OnChildEntityRemoved(one);
                    continue;
                }

                i++;
            }
        }

        private void OnChildEntityRemoved(Entity child)
        {   
        }

        // 立刻调用所属和子节点所属Component.Distroy
        public void Destroy()
        {
            if (_destroyed)
                return;
           
            destroyComponents();
            destroyChilds();
            _destroyed = true;

            // 等待父节点.Update再从节点树中移除
            // 立刻从管理器移除
            _env.RemoveEntity(this);
        }

        private void destroyComponents()
        {
            for (var i = 0; i < _comps.Count; i++)
            {
                var one = _comps[i];
                if (one.destroyed)
                {
                    continue;
                }

                one.DoDestroy();
            }
        }

        private void destroyChilds()
        {
            for (var i = 0; i < _childs.Count; i++)
            {
                var one = _childs[i];
                if (one.destroyed)
                {
                    continue;
                }

                one.Destroy();
            }
        }

        public Entity parent
        {
            get { return _parent; }
            set { SetParent(value); }
        }

        // 注: 如果在updateChilds循环中立刻调用RemoveAt，可能导致
        // 某个Entity丢掉一次Update
        // 目前都在子节点更新完毕后，再移除不属于自己的子节点
        // 极端情况:
        // 在某节点Update之后，被挂到后面某个节点上面，那么可能
        // 再被Update一次
        // (非componentRunBatch模式)
        public void SetParent(Entity parent)
        {
            if (_parent == parent)
                return;
            _parent = parent;
            _parent.addChild(this);
        }

        private void addChild(Entity child)
        {
            _childs.Add(child);
        }
    }
}
