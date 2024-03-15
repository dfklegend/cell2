using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    // 1耐 = 10血
    public class HPMaxTransformer : IAttrTransformer
    {
        public string GetAttrName()
        {
            return AttrDefine.HPMax;
        }

        public int GetLevel()
        {
            return 2;
        }

        public void Transform(IAttrOwner owner)
        {
            owner.AddTransformAttr(AttrDefine.HPMax, owner.GetAttrs().GetAttr(AttrDefine.Stamina).final * 10f);
            //var attrs = owner.GetAttrs();
            //attrs.GetAttr(AttrDefine.HPMax).Transformed.baseValue +=
            //    attrs.GetAttr(AttrDefine.Stamina).final * 10f;
        }

        public static HPMaxTransformer It = new HPMaxTransformer();
    }
}// namespace Phoenix
