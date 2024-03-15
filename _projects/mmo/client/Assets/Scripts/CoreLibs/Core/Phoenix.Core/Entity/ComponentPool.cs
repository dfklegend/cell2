using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Phoenix.Core
{    
    public class ComponentPool
    {
        const int ALLOC_STEP = 100;
        private Type _type;
        private int _activeNum = 0;

        // 通过
        private Queue<BaseComponent> _comps = new Queue<BaseComponent>();

        public ComponentPool(Type type)
        {
            _type = type;
        }

        public T Alloc<T>()
            where T: BaseComponent
        {
            _activeNum++;
            if (_comps.Count > 0)
            {
                return _comps.Dequeue() as T;
            }

            extendPool();
            return _comps.Dequeue() as T;
        }

        public BaseComponent alloc()
        {
            return Activator.CreateInstance(_type) as BaseComponent;
        }

        public void Free(BaseComponent c)
        {
            _activeNum--;
            _comps.Enqueue(c);
        }

        private int analysisStep(int activeNum)
        {
            if (activeNum < 5)
                return 1;
            if (activeNum < 10)
                return 10;
            if (activeNum < 100)
                return 30;
            return ALLOC_STEP;
        }

        private void extendPool()
        {
            int step = analysisStep(_activeNum);

            for(var i = 0; i < step; i ++)
            {
                _comps.Enqueue(alloc());
            }
        }        
    }    

    public class ComponentPoolMgr
    {
        private Dictionary<Type, ComponentPool> _pools = new Dictionary<Type, ComponentPool>();

        public T Alloc<T>()
            where T : BaseComponent
        {
            ComponentPool pool;
            var type = typeof(T);
            if (_pools.TryGetValue(type, out pool))
            {
                return pool.Alloc<T>();
            }

            pool = addPool(type);
            return pool.Alloc<T>();
        }

        public void Free(BaseComponent c)
        {
            var type = c.GetType();
            ComponentPool pool;
            if(_pools.TryGetValue(type, out pool))
            {
                pool.Free(c);
            }
            else
            {
                pool = addPool(type);
                pool.Free(c);
            }
        }

        private ComponentPool addPool(Type type)
        {
            var pool = new ComponentPool(type);
            _pools[type] = pool;
            return pool;
        }
    }    
}
