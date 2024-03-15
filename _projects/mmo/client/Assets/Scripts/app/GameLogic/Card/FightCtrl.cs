using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.Card
{	    
    // 提供公共的战斗场景，角色接口
    public class FightCtrl : Singleton<FightCtrl>
    {
        // shortcuts
        private Entity.Entity _playerEntity;        
        private Character _playerChar;
        public Character player => _playerChar;
        public int entityDown;

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
            _playerChar = null;
            _playerEntity = null;
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

        public void SetEntityDown(int id)
        {
            this.entityDown = id;
        }

        public int GetEntityDown()
        {
            return this.entityDown;
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

        public void DestroyEntity(int id)
        {
            var world = WorldMgr.It.GetWorld(1);
            var e = world.GetEntity(id);
            if (e == null)
                return;
            world.DestroyEntity(id);
        }

        public void PreparePlayer(int id)
        {
            this._playerEntity = GetEntity(id);            
            this._playerChar = GetChar(id);
        }

        public int GetPlayerSide()
        {
            if (this._playerChar == null)
                return -1;
            return this._playerChar.side;
        }
    }
    
} // namespace Phoenix
