using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    public class EventDefine
    {

        public const int EventBegin = 100;
        // 初始化单位
        public const int InitUnit       = EventBegin + 1;
        // 攻击
        public const int Attack         = EventBegin + 2;
        public const int HPChanged      = EventBegin + 3;
        public const int StartSkill     = EventBegin + 4;
        public const int InitFight      = EventBegin + 5;
           
    }
}// namespace Phoenix
