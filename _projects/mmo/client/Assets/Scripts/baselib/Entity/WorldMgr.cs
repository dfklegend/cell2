using System.Linq;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;

namespace Phoenix.Entity
{
    public class OneWorld
    {
        public EntityWorld world;
        public Transform node;
    }

    public class WorldMgr : Singleton<WorldMgr>
    {
        private Transform _rootNode;
        private Dictionary<int, OneWorld> _worlds =
            new Dictionary<int, OneWorld>();

        private int _nextID = 1;

        public void Init(Transform root)
        {
            _rootNode = root;
        }

        public EntityWorld CreateWorld()
        {
            OneWorld w = new OneWorld();
            int ID = allocID();
            w.node = createNode(ID);

            var eworld = new EntityWorld();
            eworld.Init(ID, w.node);

            w.world = eworld;

            _worlds[ID] = w;
            return eworld;
        }

        public EntityWorld GetWorld(int ID)
        {
            OneWorld w;
            if (!_worlds.TryGetValue(ID, out w))
                return null;
            return w.world;
        }

        private int allocID()
        {
            return _nextID ++;
        }

        private Transform createNode(int ID)
        {
            var go = new GameObject();

            go.name = $"world{ID}";
            go.transform.parent = _rootNode;
            go.transform.localPosition = Vector3.zero;

            return go.transform;
        }

        public void Update()
        {
            OneWorld w;
            var keys = _worlds.Keys.ToList();
            for(var i = 0; i < keys.Count; i ++)
            {
                if (!_worlds.TryGetValue(keys[i], out w))
                    continue;
                if (w.world.IsOver())
                {
                    destroyWorldNode(w);
                    _worlds.Remove(keys[i]);
                }
                else
                    w.world.Update();
            }
        }

        private void destroyWorldNode(OneWorld w)
        {
            if (!w.node)
                return;
            GameObject.Destroy(w.node.gameObject);
        }
    }
} // namespace Phoenix
