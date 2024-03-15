using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.Card
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
        public Entity.Entity CreatePlayer(EntityWorld world,
            MonsterCreateInfo info)
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
            c.SetId(e.GetEntityID());
            c.SetSide(info.side);
            //c.InitMonster(info.monsterId, info.level);        
            u.SetCharacter(c);

            //c.Dump();
            // create model
            UnitModelMgr.It.CreateModel(u);            
            return e;
        }

        public Entity.Entity CreateMonster(EntityWorld world,
            MonsterCreateInfo info)
        {           
            var e = world.CreateEntity() as Entity.Entity;
            if (e == null)
                return null;
            var cfg = MonsterDataMgr.It.GetItem(info.monsterId);
            if (cfg == null)
                return null;
            Log.LogCenter.Default.Debug("CreateCreature");
            var u = new CharacterUnit();
            u.SetUnitType((int)eUnitType.Creature);
            e.BindLogicUnit(u);

            // character
            Character c = new Character();            
            c.SetId(e.GetEntityID());
            c.SetSide(info.side);
            //c.InitMonster(info.monsterId, info.level);        
            u.SetCharacter(c);

            //c.Dump();
            // create model
            var m = UnitModelMgr.It.CreateModel(u);
            m.SetIcon(CardUtils.MakeIconPath(cfg.icon));
            return e;
        }        
    }

} // namespace Phoenix
