using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.Card
{
    public class EventDefine
    {

        public const int EventBegin = 20000;


        // ---- network [200-300]
        public const int NetworkBegin   = EventBegin + 100;
        
        public const int LoginSucc      = NetworkBegin;
        public const int ConnectFailed  = NetworkBegin + 1;
        public const int ConnectionBroken  = NetworkBegin + 2;

        public const int NetworkEnd     = NetworkBegin + 99;
        // ----

        // ---- logic
        public const int LogicBegin = EventBegin + 200;

        public const int CharInfo = LogicBegin;
        public const int BattleLog = LogicBegin + 1;
        // 初始化单位
        public const int InitUnit = LogicBegin + 2;
        // 更新场景信息
        public const int RefreshSceneInfo = LogicBegin + 3;


        // ----

        // ---- fight
        public const int FightBegin = EventBegin + 300;

        public const int TestSnapshot = FightBegin;
        public const int MoveTo = FightBegin + 1;
        public const int UnitLeave = FightBegin + 2;
        public const int Attack = FightBegin + 3;
        public const int UnitRelive = FightBegin + 4;
        public const int ServerStartSkill = FightBegin + 5;
        public const int ServerSkillHit = FightBegin + 6;    
        public const int HPChanged = FightBegin + 7;
        public const int StartSkill = FightBegin + 8;
        public const int InitFight = FightBegin + 9;
        public const int UnitAttrsChanged = FightBegin + 10;
        public const int MPChanged = FightBegin + 11;
        public const int ServerSkillBroken = FightBegin + 12;
        public const int SkillBroken = FightBegin + 13;
        public const int StopMove = FightBegin + 14;
        public const int ExitSnapshot = FightBegin + 15;
        public const int MonsterSnapshot = FightBegin + 16;

        // ----
        // ---- systems
        public const int SystemsBegin = EventBegin + 500;
        public const int RefreshServerCards = SystemsBegin;
        public const int RefreshCards = SystemsBegin + 1;
        // ----
    }
}// namespace Phoenix
