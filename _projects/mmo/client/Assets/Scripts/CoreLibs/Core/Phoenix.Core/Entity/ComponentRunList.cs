using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Phoenix.Core
{
    // 当前对应的控件执行列表
    public class ComponentRunList
    {
        private Type _type;
        public Type type { get{ return _type;} }
        // executeOrder和needUpdate不能动态修改
        public int executeOrder = 0;
        // 优化，纯数据component，没有定义OnUpdate的不刷新
        public bool needUpdate = true;

        private List<BaseComponent> _comps = new List<BaseComponent>();

        public ComponentRunList(Type type)
        {
            _type = type;
        }

        public void Add(BaseComponent c)
        {
            if (!needUpdate)
                return;
            _comps.Add(c);
        }

        public void Update()
        {
            if (!needUpdate)
                return;
            for (var i = 0; i < _comps.Count; i ++)
            {
                var one = _comps[i];
                if (!one.destroyed)
                    one.Update();
            }
        }

        // 统一调用
        public void RecycleDestroyed(ComponentPoolMgr pools)
        {
            if (!needUpdate)
                return;
            for (var i = 0; i < _comps.Count;)
            {
                var one = _comps[i];
                if (one.destroyed)
                {
                    one.Recycle();
                    pools.Free(one);
                    _comps.RemoveAt(i);
                    continue;
                }

                i++;
            }
        }
    }
    
    public class ComponentRunMgr
    {
        private Dictionary<Type, ComponentRunList> _lists = new Dictionary<Type, ComponentRunList>();
        private List<ComponentRunList> _values;
        private bool _valuesDirt = false;

        public void Add(BaseComponent c)
        {
            ComponentRunList runList;
            var type = c.GetType();
            if (_lists.TryGetValue(type, out runList))
            {
                runList.Add(c);
                return;
            }

            runList = addList(type);

            runList.needUpdate = componentNeedUpdate(c);
            runList.executeOrder = c.GetExecuteOrder();

            runList.Add(c);
            _valuesDirt = true;
        }        

        private bool componentNeedUpdate(BaseComponent c)
        {
            if (c.ForceNoUpdate())
                return false;
            var type = c.GetType();
            // 有OnUpdate，或者OnAwake函数定义
            if (Utils.ReflectUtil.HasOverrideFunc(type, typeof(BaseComponent), "OnUpdate")
                || Utils.ReflectUtil.HasOverrideFunc(type, typeof(BaseComponent), "OnAwake") )
                return true;
            return false;
        }

        private ComponentRunList addList(Type type)
        {
            var list = new ComponentRunList(type);
            _lists[type] = list;
            return list;
        }

        public void Update()
        {
            // Update中可能增加            
            if(_valuesDirt)
            {
                _values = _lists.Values.ToList();
                // 
                _values.Sort((a, b) => { return a.executeOrder - b.executeOrder; });
                _valuesDirt = false;
            }
            
            foreach ( var one in _values)
            {
                one.Update();
            }
        }

        public void RecycleDestroyed(ComponentPoolMgr pools)
        {   
            foreach (var one in _values)
            {
                one.RecycleDestroyed(pools);
            }
        }
    }
}
