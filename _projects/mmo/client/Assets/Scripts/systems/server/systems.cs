using System;
using System.Collections.Generic;
using Phoenix.Core;
using Phoenix.Utils;

namespace Phoenix.Game
{
    public partial class Systems : Singleton<Systems>
    {
        private Dictionary<string, BaseSystem> _systems = new Dictionary<string, BaseSystem>();

        public void OnServerCmd(string system, string cmd, byte[] argsBuf)
        {
            BaseSystem sys;
            if(!_systems.TryGetValue(system, out sys) )        
            {
                Log.LogCenter.Default.Debug($"Systems has not system: {system}");    
                return;
            }
            sys.OnCmd(cmd, argsBuf);
        }

        public void AddSystem(BaseSystem system)
        {
            _systems[system.GetName()] = system;
        }

        public BaseSystem GetSystem(string system)
        {
            return _systems[system];
        }

        public T GetSystem<T>(string system)
            where T: class
        {
            var s = GetSystem(system);
            if (s == null)
                return null;
            return s as T;
        }

        public void Reset()
        {
            foreach(KeyValuePair<string, BaseSystem> kv in _systems)
            {
                kv.Value.Reset();
            }
        }
    }
}
