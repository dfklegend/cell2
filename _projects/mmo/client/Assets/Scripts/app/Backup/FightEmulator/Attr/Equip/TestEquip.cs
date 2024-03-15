using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    /*
     * 测试装备
     *    +5 hp
     *    +3 attack
     *    +2 str
     */
    public class TestEquip : IAttrEquipable
    {
        public void Equip(IAttrOwner owner)
        {
            applyEquip(owner, 1f);
        }

        public void Unequip(IAttrOwner owner)
        {
            applyEquip(owner, -1f);
        }

        private void applyEquip(IAttrOwner owner, float factor)
        {
            var attrs = owner.GetAttrs();
            attrs.GetAttr(AttrDefine.HPMax).Base.baseValue += 5*factor;
            attrs.GetAttr(AttrDefine.MeleePower).Base.baseValue += 3*factor;
            attrs.GetAttr(AttrDefine.Strength).Base.baseValue += 2*factor;
        }
    }

    /*
     * 测试装备
     *    hpMax +10%
     *    attack +10%
     *    str +10%
     *    str +10
     */
    public class TestEquip1 : IAttrEquipable
    {
        public void Equip(IAttrOwner owner)
        {
            applyEquip(owner, 1f);
        }

        public void Unequip(IAttrOwner owner)
        {
            applyEquip(owner, -1f);
        }

        private void applyEquip(IAttrOwner owner, float factor)
        {
            var attrs = owner.GetAttrs();
            AttrUtil.AddAllPercent(attrs.GetAttr(AttrDefine.HPMax), 0.1f * factor);
            AttrUtil.AddAllPercent(attrs.GetAttr(AttrDefine.MeleePower), 0.1f * factor);
            AttrUtil.AddAllPercent(attrs.GetAttr(AttrDefine.Strength), 0.1f * factor);
            attrs.GetAttr(AttrDefine.MeleePower).Base.baseValue += 10 * factor;
            attrs.GetAttr(AttrDefine.Strength).Base.baseValue += 10 * factor;
        }
    }
}// namespace Phoenix
