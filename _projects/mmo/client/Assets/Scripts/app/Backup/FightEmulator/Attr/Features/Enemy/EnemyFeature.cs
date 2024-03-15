using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{   
    public class EnemyFeature : IAttrFeature
    {
        public static EnemyFeature It = new EnemyFeature();

        public void ApplyInitAttrs(IAttrOwner owner)
        {   
        }

        public void ApplyLevelupAttrs(IAttrOwner owner, int offset)
        {
            string monsterId = owner.GetMonsterId();
            MonsterData ed = MonsterDataMgr.It.GetItem(monsterId);
            if (ed == null)
                return;

            int off = offset + 1 - ed.baseLevel;
            var attrs = owner.GetAttrs();
            var weapon = owner.GetEnemyWeapon() as EnemyWeapon;

            weapon.Reset();
            weapon.SetSpeed(ed.weaponSpeed);
            weapon.OffMinDmg(ed.minDmg + ed.minDmgLv * off);
            weapon.OffMaxDmg(ed.maxDmg + ed.maxDmgLv * off);

            attrs.GetAttr(AttrDefine.Armor).Reset();
            attrs.GetAttr(AttrDefine.Armor).Base.baseValue += ed.armor + ed.armorLv*off;
            attrs.GetAttr(AttrDefine.HPMax).Reset();
            attrs.GetAttr(AttrDefine.HPMax).Base.baseValue += ed.hp + ed.hpLv*off;
        }
    }
}// namespace Phoenix
