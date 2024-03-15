using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.Card
{ 
    // 提供一个serverId : entityId映射关系
    public class ServerCtrl : Singleton<ServerCtrl>
    {
        private Dictionary<int, int> _map = new Dictionary<int, int>();

        public void Add(int serverId, int entityId)
        {
            _map[serverId] = entityId;
        }

        public int GetEntityId(int serverId) 
        {            
            int entityId;
            if (!_map.TryGetValue(serverId, out entityId))
                return -1;
            return entityId;
        }

        public void Remove(int serverId)
        {
            _map.Remove(serverId);
        }

        public void Clear()
        {
            _map.Clear();
        }
    }

} // namespace Phoenix
