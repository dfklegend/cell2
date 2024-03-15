using System;
using System.Collections.Generic;
using System.IO;
using System.Text;
using Phoenix.Utils;

namespace Phoenix.Core
{
    public class SceneMgr : BaseComponent
    {
        RunEnv _env;
        private Entity _sceneRoot;
        Dictionary<ulong, Entity> _scenes =
            new Dictionary<ulong, Entity>();
        List<ulong> _toRemove = new List<ulong>();

        public static SceneMgr Create(RunEnv env)
        {
            Entity root = createScenesEntity(env);
            var mgr = root.AddComponent<SceneMgr>();
            mgr._env = env;
            mgr._sceneRoot = root;
            return mgr;
        }

        private static Entity createScenesEntity(RunEnv env)
        {
            return env.CreateEntity(null);
        }

        public Entity CreateScene()            
        {
            Entity e = _env.CreateEntity(_sceneRoot);            
            _scenes[e.id] = e;
            return e;
        }

        public Entity GetScene(ulong id)
        {
            Entity scene;
            if (_scenes.TryGetValue(id, out scene))
                return scene;
            return null;
        }

        public void DestroyScene(Entity scene)
        {
            scene.Destroy();
        }

        public override void OnUpdate()
        {
            base.OnUpdate();
            removeDestroyedScenes();
        }        

        private void removeDestroyedScenes()
        {
            foreach(var key in _scenes.Keys)
            {
                if (_scenes[key].destroyed)
                    _toRemove.Add(key);
            }

            if (_toRemove.Count == 0)
                return;
            for(var i = 0; i < _toRemove.Count; i ++)
            {
                _scenes.Remove(_toRemove[i]);
            }
            _toRemove.Clear();
        }
    }
}