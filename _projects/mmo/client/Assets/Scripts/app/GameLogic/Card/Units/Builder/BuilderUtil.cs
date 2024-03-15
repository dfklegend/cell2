using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.Card
{	
    public static class BuilderUtil
    {
        public static int CreateMonster(EntityWorld world,
            string cfgId, int level, int side, float x, float z)
        {
            var info = new MonsterCreateInfo();
            info.monsterId = cfgId;
            info.level = level;
            info.side = side;

            var entity = CharBuilder.It.CreateMonster(world, info);
            if (entity == null)
                return -1;
            var c = GetCharFromEntity(entity);
            //c.Dump();

            c.pos = new Vector3(x, 0, z);
            UnitModelMgr.It.UpdateModelPos(entity.ID);
            return entity.ID;
        }

        public static int CreatePlayer(EntityWorld world,
            int level, int side, float x, float z)
        {
            var info = new MonsterCreateInfo();            
            info.level = level;
            info.side = side;

            var entity = CharBuilder.It.CreatePlayer(world, info);
            if (entity == null)
                return -1;
            var c = GetCharFromEntity(entity);
            //c.Dump();

            c.pos = new Vector3(x, 0, z);
            UnitModelMgr.It.UpdateModelPos(entity.ID);
            return entity.ID;
        }

        public static Character GetCharFromEntity(Entity.Entity e)
        {
            if (e == null)
                return null;
            var u = e.unit as CharacterUnit;
            return u.character;
        }
    }
} // namespace Phoenix
