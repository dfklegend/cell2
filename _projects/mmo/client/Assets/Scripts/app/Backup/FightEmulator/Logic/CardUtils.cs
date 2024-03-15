using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;



namespace Phoenix.Game.FightEmulator
{
    public static class CardUtils
    {
        // 设置血量            
        public static void CharSetHP(Character c, int hp)
        {
            c.attrs.GetAttr(AttrDefine.HP).Base.baseValue = hp;
        }

        // 同步血量变化
        public static void SyncHPBar(Character c)
        {
            HEventUtil.Dispatch(GlobalEvents.It.events,
               new HEventHPChanged(c));
        }

        public static void SyncInfo(Character c)
        {
            HEventUtil.Dispatch(GlobalEvents.It.events,
               new HEventInitUnit(c.id));
        }

        // 播放攻击动画
        public static void PlayAttack(Character src, Character tar)
        {
            HEventUtil.Dispatch(GlobalEvents.It.events,
                new HEventStartSkill(src, "普攻", tar.id));
        }

        public static void PlaySimpleHit(Character src, Character tar, int dmg, bool critical)
        {
            var result = new FormulaResult();
            result.data.Dmg = dmg;
            if (dmg > 0)
                result.data.hit = true;
            else
                result.data.hit = false;
            result.data.critical = critical;

            HEventUtil.Dispatch(GlobalEvents.It.events,
                new HEventAttack(src, tar, result));
        }
    }
        
} // namespace Phoenix
