using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    // 攻击 = 力量*2    
    public class TestTransformer : IAttrTransformer
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
            var newV = attrs.GetAttr(AttrDefine.Strength).final * 2.0f;
            owner.AddTransformAttr(AttrDefine.MeleePower, newV);
        }
    }
}// namespace Phoenix
