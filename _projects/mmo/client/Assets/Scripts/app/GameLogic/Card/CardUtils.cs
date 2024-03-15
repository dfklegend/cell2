using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;



namespace Phoenix.Game.Card
{
    public static class CardUtils
    {
        // 设置血量            
        public static void CharSetHP(Character c, int hp)
        {
            c.SetHP(hp);
        }

        // 同步血量变化
        public static void SyncHPBar(Character c)
        {
            HEventUtil.Dispatch(GlobalEvents.It.events,
               new HEventHPChanged(c));
        }

        public static void SyncMPBar(Character c)
        {
            HEventUtil.Dispatch(GlobalEvents.It.events,
               new HEventMPChanged(c));
        }

        public static void SyncInfo(Character c)
        {
            HEventUtil.Dispatch(GlobalEvents.It.events,
               new HEventInitUnit(c.id));
        }

        // 播放攻击动画
        public static void PlayAttack(string skillId, Character src, Character tar)
        {
            HEventUtil.Dispatch(GlobalEvents.It.events,
                new HEventStartSkill(src, skillId, tar.id));
        }

        public static void PlaySimpleHit(string skillId, Character src, Character tar, int dmg, bool critical)
        {
            var result = new FormulaResult();
            result.data.Dmg = dmg;
            if (dmg > 0)
                result.data.hit = true;
            else
                result.data.hit = false;
            result.data.critical = critical;

            HEventUtil.Dispatch(GlobalEvents.It.events,
                new HEventAttack(skillId, src, tar, result));
        }

        public static void DoSkillBroken(Character c, string skillId)
        {
            HEventUtil.Dispatch(GlobalEvents.It.events,
               new HEventSkillBroken(c, skillId));
        }

        public static void loadAllPrefabs()
        {
            ItemPrefabFactory.It.Init();
            ItemListStyleFactory.It.Init();
            ItemIconStyleFactory.It.Init();
        }

        public static void loadAllTables()
        {
            TestDataMgr.It.Load();
            ItemDataMgr.It.Load();
            EquipDataMgr.It.Load();
            SkillDataMgr.It.Load();
            MonsterDataMgr.It.Load();
            MonsterTemplateDataMgr.It.Load();
        }

        public static string MakeIconPath(string name)
        {
            return $"icons/heros/{name}";
        }
    }
        
} // namespace Phoenix
