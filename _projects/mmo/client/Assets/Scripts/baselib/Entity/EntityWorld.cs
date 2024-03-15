using System.Linq;
using System.Collections.Generic;
using UnityEngine;


namespace Phoenix.Entity
{	    
    public class EntityWorld : IEntityWorld
    {
        private bool _over = false;
        private Transform _rootNode;
        public Transform rootNode { get { return _rootNode; } }
        private int _nextID = 1;
        private Dictionary<int, IEntity> _entities =
            new Dictionary<int, IEntity>();

        private int _worldID = 0;
        public int worldID { get { return _worldID; } }

        public void Init(int worldID, Transform root)
        {
            _worldID = worldID;
            if (!root)
                return;
            _rootNode = root;
        }

        private void resetIDAlloc()
        {
            _nextID = 1;
        }

        public IEntity CreateEntity()
        {
            var go = new GameObject();
            var e = go.AddComponent<Entity>();

            // set id
            e.SetEntityID(allocID());
            e.SetWorldID(worldID);

            go.name = $"entity{e.GetEntityID()}";
            // attach
            go.transform.parent = _rootNode;
            go.transform.localPosition = Vector3.zero;            
            
            _entities[e.GetEntityID()] = e;
            return e;
        }

        private int allocID()
        {
            return _nextID++;
        }

        public IEntity GetEntity(int ID)
        {
            IEntity ret;
            if(!_entities.TryGetValue(ID, out ret))
                return null;
            return ret;
        }

        public void DestroyEntity(int ID)
        {
            var one = GetEntity(ID);
            if (one == null || one.IsOver())
                return;
            one.Destroy();
        }

        public void Update()
        {   
            var keys = _entities.Keys.ToList();
            IEntity e;
            for (var i = 0; i < keys.Count; i ++)
            {
                var key = keys[i];                
                if (!_entities.TryGetValue(key, out e))
                    continue;
                if(e.IsOver())
                {
                    removeEntity(e);
                }
            }
        }

        public void Visit(IVisitor visitor)
        {
            var keys = _entities.Keys.ToList();
            IEntity e;
            for (var i = 0; i < keys.Count; i++)
            {
                var key = keys[i];
                if (!_entities.TryGetValue(key, out e))
                    continue;
                visitor.Visit(e);
            }
        }

        private void removeEntity(IEntity e)
        {            
            _entities.Remove(e.GetEntityID());
        }

        public void Clear()
        {
            IEntity e;
            var keys = _entities.Keys.ToList();
            for (var i = 0; i < keys.Count; i++)
            {
                var key = keys[i];                
                if (!_entities.TryGetValue(key, out e))
                    continue;
                e.Destroy();                
            }
            _entities.Clear();
        }

        public void Reset()
        {
            Clear();
            resetIDAlloc();
        }

        public void Destroy()
        {
            if (IsOver())
                return;
            setOver();
            Clear();
        }

        public bool IsOver()
        {
            return _over;
        }

        private void setOver()
        {
            _over = true;
        }
    }
    
} // namespace Phoenix
