using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;

namespace Phoenix.Game.FightEmulator
{
    public class HEventStartSkill: HEvent<int>
    {
        public Character src;        
        public string skillId;
        public int tarId;


        public HEventStartSkill(Character src, string skillId, int tarId)
            : base(EventDefine.StartSkill)
        {
            this.src = src;            
            this.skillId = skillId;
            this.tarId = tarId;
        }
    }

    
}// namespace Phoenix
