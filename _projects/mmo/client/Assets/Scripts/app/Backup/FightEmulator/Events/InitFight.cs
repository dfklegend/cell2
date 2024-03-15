using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;


namespace Phoenix.Game.FightEmulator
{
    public class HEventInitFight: HEvent<int>
    {
        public HEventInitFight()
            : base(EventDefine.InitFight)
        {   
        }
    }
}// namespace Phoenix
