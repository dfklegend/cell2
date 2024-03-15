using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    // 力量 = 8 + 1.8*level
    // 敏捷 = 8+ 1.8*level
    // 耐力 = 5 + 1.0*level
    // 智力 = 1 + 0.4*level
    // HPMax = 100 + 20*level
    public class RogueFeature : IAttrFeature
    {
        public void ApplyInitAttrs(IAttrOwner owner)
        {
            var attrs = owner.GetAttrs();
            attrs.GetAttr(AttrDefine.Strength).Base.baseValue += 8;
            attrs.GetAttr(AttrDefine.Agility).Base.baseValue += 8;
            attrs.GetAttr(AttrDefine.Stamina).Base.baseValue += 5;
            attrs.GetAttr(AttrDefine.Intellect).Base.baseValue += 1;

            attrs.GetAttr(AttrDefine.HPMax).Base.baseValue += 100;
        }

        public void ApplyLevelupAttrs(IAttrOwner owner, int offset)
        {
            var attrs = owner.GetAttrs();            

            owner.AddAttr(AttrDefine.Strength, false, eElementType.Base, 1.8f * offset);
            owner.AddAttr(AttrDefine.Agility, false, eElementType.Base, 1.8f * offset);
            owner.AddAttr(AttrDefine.Stamina, false, eElementType.Base, 1.0f * offset);
            owner.AddAttr(AttrDefine.Intellect, false, eElementType.Base, 0.4f * offset);
            owner.AddAttr(AttrDefine.HPMax, false, eElementType.Base, 20 * offset);            
        }
    }
}// namespace Phoenix
