using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;


namespace Phoenix.Game.FightEmulator
{
    public class HEventAttack: HEvent<int>
    {
        public Character src;
        public Character tar;
        public FormulaResult result;


        public HEventAttack(Character src, Character tar, FormulaResult result)
            : base(EventDefine.Attack)
        {
            this.src = src;
            this.tar = tar;
            this.result = result;
        }
    }
}// namespace Phoenix
