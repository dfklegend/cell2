using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;

namespace Phoenix.Core
{
    public class EntityMgr
    {
        private Dictionary<ulong, Entity> _entities = new Dictionary<ulong, Entity>();

        public void AddEntity(Entity e)
        {
            _entities[e.id] = e;
        }

        public void RemoveEntity(Entity e)
        {
            _entities.Remove(e.id);
        }

        public Entity GetEntity(ulong id)
        {
            Entity e;
            if (_entities.TryGetValue(id, out e))
                return e;
            return null;
        }
    }
}
