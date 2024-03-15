using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    // 力量 = 10 + 1.5*level
    // HPMax = 100 + 10*level
    public class TestFeature : IAttrFeature
    {
        public void ApplyInitAttrs(IAttrOwner owner)
        {
            var attrs = owner.GetAttrs();
            attrs.GetAttr(AttrDefine.Strength).Base.baseValue += 10;
            attrs.GetAttr(AttrDefine.HPMax).Base.baseValue += 100;
        }

        public void ApplyLevelupAttrs(IAttrOwner owner, int offset)
        {
            var attrs = owner.GetAttrs();
            attrs.GetAttr(AttrDefine.Strength).Base.baseValue += 1.5f*offset;
            attrs.GetAttr(AttrDefine.HPMax).Base.baseValue += 10*offset;
        }
    }
}// namespace Phoenix
