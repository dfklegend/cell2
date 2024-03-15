using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;

namespace Phoenix.Game.FightEmulator
{
    public class HEventHPChanged: HEvent<int>
    {
        public Character src;

        public HEventHPChanged(Character src)
            : base(EventDefine.HPChanged)
        {
            this.src = src;           
        }
    }
}// namespace Phoenix
