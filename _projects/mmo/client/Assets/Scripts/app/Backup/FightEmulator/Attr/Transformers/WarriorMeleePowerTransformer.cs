using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    // 攻击 = level*3 + 力量*2 - 20    
    public class WarriorMeleePowerTransformer : IAttrTransformer
    {
        public string GetAttrName()
        {
            return AttrDefine.MeleePower;
        }

        public int GetLevel()
        {
            return 2;
        }

        public void Transform(IAttrOwner owner)
        {
            var attrs = owner.GetAttrs();
            var newV = owner.GetLevel()*3f +
                attrs.GetAttr(AttrDefine.Strength).final * 2.0f - 20;
            owner.AddTransformAttr(AttrDefine.MeleePower, newV);
        }

        public static WarriorMeleePowerTransformer It = new WarriorMeleePowerTransformer();
    }
}// namespace Phoenix
