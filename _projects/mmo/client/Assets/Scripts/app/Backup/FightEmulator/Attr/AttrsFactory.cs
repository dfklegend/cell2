using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    public static class AttrsFactory
    {   
        private static IAttrFinalClamper clampHPMax = new HPMaxClamper();
        
        // 后续组织成pool
        // 计算中可能存在需要大量临时Attrs的情况
        public static Attrs Alloc()
        {
            var one = new Attrs();
            AttrDefine.InitAttrs(one);

            // 针对属性，设置IElementValueClamper
            one.GetAttr("HPMax").SetClamper(clampHPMax);            

            return one;
        }

        public static void Free(Attrs one)
        {
            one.Reset();
        }
    }
}// namespace Phoenix
