using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{	    
    // 提供公共的战斗场景，角色接口
    public class FightCtrl : Singleton<FightCtrl>
    {
        // shortcuts
        public Entity.Entity playerEntity;
        public PlayerComponent player;
        public Character playerChar;

        public EntityWorld GetWorld()
        {
            return WorldMgr.It.GetWorld(1);
        }

        public void Prepare()
        {
            InitFight();
        }

        public void InitFight()
        {   
            var world = WorldMgr.It.GetWorld(1);
            world.Reset();

            SearchCtrl.It.Init(world);            
        }
       
        public void Restart()
        {
            InitFight();
        }

        public void Destroy()
        {
            var world = WorldMgr.It.GetWorld(1);
            world.Reset();
        }

        public Entity.Entity GetEntity(int id)
        {
            var world = WorldMgr.It.GetWorld(1);
            var e = world.GetEntity(id);
            if (e == null)
                return null;
            return e as Entity.Entity;
        }

        public Character GetChar(int id)
        {
            var world = WorldMgr.It.GetWorld(1);
            var e = world.GetEntity(id);
            if (e == null)
                return null;
            Entity.Entity entity = e as Entity.Entity;
            if (entity == null)
                return null;
            var u = entity.unit as CharacterUnit;
            return u.character;
        }

        public PlayerComponent GetPlayer(int id)
        {
            var world = WorldMgr.It.GetWorld(1);
            var e = world.GetEntity(id);
            if (e == null)
                return null;
            Entity.Entity entity = e as Entity.Entity;
            return entity.GetComponent<PlayerComponent>();
        }

        public void PreparePlayer(int id)
        {
            this.playerEntity = GetEntity(id);
            this.player = GetPlayer(id);
            this.playerChar = GetChar(id);
        }

        public void DestroyEntity(int id)
        {
            var world = WorldMgr.It.GetWorld(1);
            var e = world.GetEntity(id);
            if (e == null)
                return;
            world.DestroyEntity(id);
        }
    }
    
} // namespace Phoenix
