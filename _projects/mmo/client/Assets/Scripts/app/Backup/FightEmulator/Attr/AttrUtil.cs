using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    public static class AttrUtil
    {   
        public static void AddAllPercent(Attr attr, float v)
        {
            attr.VisitElements((element) => 
            {
                element.percent += v;
            });
        }
    }
}// namespace Phoenix
