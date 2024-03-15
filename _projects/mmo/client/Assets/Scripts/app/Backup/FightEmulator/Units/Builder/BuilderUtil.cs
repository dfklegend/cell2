using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
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
            var c = GetCharFronEntity(entity);
            //c.Dump();

            c.pos = new Vector3(x, 0, z);
            UnitModelMgr.It.UpdateModelPos(entity.ID);
            return entity.ID;
        }

        public static Character GetCharFronEntity(Entity.Entity e)
        {            
            if (e == null)
                return null;           
            var u = e.unit as CharacterUnit;
            return u.character;
        }

        public static void CreateWarrior(EntityWorld world,
            int level, int side)
        {   
            RoleCreateInfo info = new RoleCreateInfo();
            info.name = "warrior";
            info.level = level;
            info.classType = (int)eClass.Warrior;
            info.side = side;
            var entity = CharBuilder.It.CreateRole(world, info);
            if (entity == null)
                return;
            var warrior = GetCharFronEntity(entity);

            EquipItem equip1 = null;            
            equip1 = EquipUtil.CreateEquip("短剑_1");           

            EquipItem equip2 = null;            
            equip2 = EquipUtil.CreateEquip("盾牌_1");            

            if (equip1 != null)
                warrior.Equip(eEquipSlot.MainHand, -1, equip1);
            if (equip2 != null)
                warrior.Equip(eEquipSlot.OffHand, -1, equip2);
            warrior.pos = new Vector3(0, 0, 0);
            UnitModelMgr.It.UpdateModelPos(entity.ID);

            var player = entity.GetComponent<PlayerComponent>();
            var bag = player.bags.GetBag((int)BagSystem.eBagType.BagItem);

            bag.AddItem(ItemSystem.ItemBuilder.CreateEquip("短剑_1"));
            bag.AddItem(ItemSystem.ItemBuilder.CreateEquip("匕首_1"));
            bag.AddItem(ItemSystem.ItemBuilder.CreateEquip("木棍_1"));
            bag.AddItem(ItemSystem.ItemBuilder.CreateEquip("盾牌_1"));

            warrior.SetBags(player.bags);

            //warrior.Dump();            
        }
    }
} // namespace Phoenix
