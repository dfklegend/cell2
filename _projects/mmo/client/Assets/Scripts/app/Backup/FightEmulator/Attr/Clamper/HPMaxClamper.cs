using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;

namespace Phoenix.Game.FightEmulator
{
    public class HPMaxClamper : IAttrFinalClamper
    {        
        public float ClampFinal(float v)
        {
            // 最小1
            return Math.Max(v, 1);
        }
    }
}// namespace Phoenix
