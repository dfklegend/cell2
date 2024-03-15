using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{	
    public class RoleCreateInfo
    {
        public string name = "";
        public int level = 1;
        public int classType = 0;
        public int side = 0;
    }

    public class MonsterCreateInfo
    {
        public string monsterId;
        public int level = 1;
        public int side = 0;
    }

    public class CharBuilder : Singleton<CharBuilder>
    {
        public Entity.Entity CreateRole(EntityWorld world,
            RoleCreateInfo info)
        {
            var e = CreateCreature(world, info);
            
            var player = e.AddComponent<PlayerComponent>();
            player.SetBags(BagSystem.PlayerBagsBuilder.CreateBags());            
            return e;
        }

        public Entity.Entity CreateCreature(EntityWorld world,
            RoleCreateInfo info)
        {
            var e = world.CreateEntity() as Entity.Entity;
            if (e == null)
                return null;
            Log.LogCenter.Default.Debug("CreateCreature");
            var u = new CharacterUnit();
            u.SetUnitType((int)eUnitType.Creature);
            e.BindLogicUnit(u);

            // character
            Character c = new Character();
            c.SetName(info.name);
            c.SetClass(info.classType);
            c.SetLevel(info.level);
            c.SetId(e.GetEntityID());
            c.SetSide(info.side);
            c.PrepareAttrs();
            u.SetCharacter(c);

            //c.Dump();

            // create model
            UnitModelMgr.It.CreateModel(u);
            return e;
        }

        public Entity.Entity CreateMonster(EntityWorld world,
            MonsterCreateInfo info)
        {
            var md = MonsterDataMgr.It.GetItem(info.monsterId);
            if (md == null)
                return null;
            var e = world.CreateEntity() as Entity.Entity;
            if (e == null)
                return null;
            Log.LogCenter.Default.Debug("CreateCreature");
            var u = new CharacterUnit();
            u.SetUnitType((int)eUnitType.Creature);
            e.BindLogicUnit(u);

            // character
            Character c = new Character();
            c.SetName(md.name);                        
            c.SetId(e.GetEntityID());
            c.SetSide(info.side);
            c.InitMonster(info.monsterId, info.level);
            c.guardRange = md.guardRange;
            c.attackRange = md.attackRange;
            u.SetCharacter(c);

            //c.Dump();
            // create model
            UnitModelMgr.It.CreateModel(u);
            return e;
        }
    }


    
} // namespace Phoenix
