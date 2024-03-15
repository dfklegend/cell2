using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System;
using Phoenix.Fsm;
using Phoenix.Utils;

namespace Phoenix.Game.FightEmulator
{
    public enum eAIState
    {
        Init = 0,     
        RandMove,           // 随机移动
        Attack,             // 攻击目标
    }    
}// namespace Phoenix
