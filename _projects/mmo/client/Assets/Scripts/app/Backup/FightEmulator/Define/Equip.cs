using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    // 装备部位
    public enum eEquipSlot
    {
        Head = 0,
        Neck,
        Shoulder,
        Chest,
        Back,
        Wrists, // 护腕
        Hands,  // 护手
        Waist,
        Legs,
        Feet,
        Finger1,
        Finger2,
        Trinket1,
        Trinket2,        
        MainHand,
        OffHand,
        Max
    }

    // 装备基础类型
    // 用雷检查装备部位合法性
    public enum eEquipBaseType
    {
        Head = 0,
        Neck,
        Shoulder,
        Chest,
        Back,
        Wrists, 
        Hands,  
        Waist,
        Legs,
        Feet,
        Finger,        
        Trinket,        
        SingleHand,     // 单手
        DoubleHand,     // 双手
        RangeWeapon,    // 远程武器
        Shield,         // 盾牌
    }
}// namespace Phoenix
