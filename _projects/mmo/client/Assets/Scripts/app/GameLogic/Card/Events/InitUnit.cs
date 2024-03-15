using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;


namespace Phoenix.Game.FightEmulator
{
    public class HEventInitUnit: HEvent<int>
    {
        public int unitId = -1;


        public HEventInitUnit(int id)
            : base(EventDefine.InitUnit)
        {
            unitId = id;
        }
    }
}// namespace Phoenix
