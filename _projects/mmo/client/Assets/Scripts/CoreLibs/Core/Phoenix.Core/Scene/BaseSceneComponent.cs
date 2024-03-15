using System;
using System.Collections.Generic;
using System.IO;
using System.Text;
using Phoenix.Utils;

namespace Phoenix.Core
{
    // 场景
    // 对象容器        
    public class BaseSceneComponent : BaseComponent
    {
        // 实际的实现        
        protected Dictionary<ulong, Entity> _entities =
            new Dictionary<ulong, Entity>();
        private List<ulong> _toRemove = new List<ulong>();

        protected ulong GetSceneId()
        {
            return owner.id;
        }

        public int GetEntityCount()
        {
            return _entities.Count;
        }

        public void AddEntity(Entity entity)
        {
            _entities[entity.id] = entity;
        }

        public void RemoveEntity(Entity entity)
        {
            _entities.Remove(entity.id);
        }

        public Entity CreateEntity()
        {
            return owner.env.CreateEntity(owner);
        }

        // 简单遍历        
        public void Visit<T>(Action<Entity, T> visitor, T arg0)
        {
            foreach(var key in _entities.Keys)
            {
                visitor(_entities[key], arg0);
            }
        }

        public virtual void OnEntityDestroy(Entity e)
        {
            
        }

        public override void OnUpdate()
        {
            base.OnUpdate();
            removeDestroyed();
        }

        private void removeDestroyed()
        {   
            foreach (var key in _entities.Keys)
            {
                var one = _entities[key];
                if (one.destroyed)
                    _toRemove.Add(key);
            }
            
            if (_toRemove.Count == 0)
                return;
            for(var i = 0; i < _toRemove.Count; i ++)
            {
                _entities.Remove(_toRemove[i]);
            }
            _toRemove.Clear();
        }
    }    
}