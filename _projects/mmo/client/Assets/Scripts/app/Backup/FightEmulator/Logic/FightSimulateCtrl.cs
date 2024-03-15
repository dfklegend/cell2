using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{	    
    // 模拟两个角色战斗
    public class FightSimulateCtrl : Singleton<FightSimulateCtrl>
    {
        public string roleClass = "warrior";
        public int roleLevel = 50;
        public string mainHand = "短剑_1";
        public string offHand = "短剑_1";

        public string enemyCfgId = "狗头人lv5";
        public int enemyLevel = 50;  

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
            Save();

            FightCtrl.It.InitFight();
            //initEntities();

            HEventUtil.Dispatch(GlobalEvents.It.events,
                new HEventInitFight());
        }

        void initEntities()
        {
            var world = GetWorld();

            RoleCreateInfo info = new RoleCreateInfo();
            info.name = "warrior";
            info.level = roleLevel;
            info.classType = (int)eClass.Warrior;
            info.side = 0;
            CharBuilder.It.CreateCreature(world, info);
            var warrior = GetChar(1);


            EquipItem equip1 = null;
            if (!string.IsNullOrEmpty(mainHand))
            {
                equip1 = EquipUtil.CreateEquip(mainHand);
            }

            EquipItem equip2 = null;
            if (!string.IsNullOrEmpty(offHand))
            {
                equip2 = EquipUtil.CreateEquip(offHand);
            }

            if (equip1 != null)
                warrior.Equip(eEquipSlot.MainHand, -1, equip1);
            if (equip2 != null)
                warrior.Equip(eEquipSlot.OffHand, -1, equip2);
            warrior.pos = new Vector3(0, 0, 10);

            warrior.Dump();
            addSkill(warrior);

            //createRogue(world);
            createMonster(world);
            //createRogue(world);
        }

        public void Load()
        {
            if (!UserPersistData.It.IsValid())
                return;
            PlayerEntry pe = UserPersistData.It.GetEntry<PlayerEntry>("player");
            // player
            roleLevel = pe.level;
            mainHand = pe.mainHand;
            offHand = pe.offHand;
            // enemy
            EnemyEntry ee = UserPersistData.It.GetEntry<EnemyEntry>("enemy");
            enemyCfgId = ee.cfgId;
            enemyLevel = ee.level;
        }

        public void Save()
        {
            PlayerEntry pe = UserPersistData.It.GetEntry<PlayerEntry>("player");
            // player
            pe.level = roleLevel;
            pe.mainHand = mainHand;
            pe.offHand = offHand;
            
            // enemy
            EnemyEntry ee = UserPersistData.It.GetEntry<EnemyEntry>("enemy");
            ee.cfgId = enemyCfgId;
            ee.level = enemyLevel;

            UserPersistData.It.SaveToFile();
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

        private void addSkill(Character c)
        {
            var slots = c.skillSlots;
            slots.SetSkill(0, "强力一击", 1);
            slots.SetSkill(1, "无伤害技能", 1);
            //slots.SetSkill(2, "技能3", 1);
            //slots.SetSkill(3, "技能4", 1);
        }

        private void createRogue(EntityWorld world)
        {
            var info = new RoleCreateInfo();
            info.name = "rogue";
            info.level = enemyLevel;
            info.classType = (int)eClass.Rogue;
            info.side = 1;
            var entity = CharBuilder.It.CreateCreature(world, info);
            var rogue = GetChar(entity.ID);

            var equip1 = new EquipItem();
            equip1.Init("短剑_1");
            var equip2 = new EquipItem();
            equip2.Init("木棍_1");         

            rogue.Equip(eEquipSlot.MainHand, -1, equip1);
            rogue.Equip(eEquipSlot.OffHand, -1, equip1);

            rogue.pos = new Vector3(5, 0, -10);

            rogue.Dump();
        }

        private void createMonster(EntityWorld world)
        {
            var info = new MonsterCreateInfo();
            info.monsterId = enemyCfgId;
            info.level = enemyLevel;
            info.side = 1;

            var entity = CharBuilder.It.CreateMonster(world, info);
            var c = GetChar(entity.ID);
            c.Dump();

            c.pos = new Vector3(0, 0, -10);
        }

        public Character GetChar(int id)
        {
            return FightCtrl.It.GetChar(id);
        }
    }
    
} // namespace Phoenix
