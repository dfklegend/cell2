using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    // 1敏 = 2护甲
    public class ArmorTransformer : IAttrTransformer
    {
        public string GetAttrName()
        {
            return AttrDefine.Armor;
        }

        public int GetLevel()
        {
            return 2;
        }

        public void Transform(IAttrOwner owner)
        {
            owner.AddTransformAttr(AttrDefine.Armor, owner.GetAttrs().GetAttr(AttrDefine.Agility).final * 2f);
            //var attrs = owner.GetAttrs();
            //attrs.GetAttr(AttrDefine.Armor).Transformed.baseValue +=
            //    attrs.GetAttr(AttrDefine.Agility).final * 2f;
        }

        public static ArmorTransformer It = new ArmorTransformer();
    }
}// namespace Phoenix
