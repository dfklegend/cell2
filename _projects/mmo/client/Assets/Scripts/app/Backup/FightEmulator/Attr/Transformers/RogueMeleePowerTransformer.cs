using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    // 攻击 = level*2 + 力量*1 + 敏捷*1 - 20
    public class RogueMeleePowerTransformer : IAttrTransformer
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
            var newV =
                owner.GetLevel() * 3f +
                attrs.GetAttr(AttrDefine.Strength).final * 1.0f +
                attrs.GetAttr(AttrDefine.Agility).final * 1.0f - 20;
            owner.AddTransformAttr(AttrDefine.MeleePower, newV);
        }

        public static RogueMeleePowerTransformer It = new RogueMeleePowerTransformer();
    }
}// namespace Phoenix
