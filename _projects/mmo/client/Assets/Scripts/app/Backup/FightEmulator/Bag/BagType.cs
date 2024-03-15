using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;

namespace Phoenix.Game.FightEmulator.BagSystem
{	
    public enum eBagType
    {
        BagTemp = 0,        // 临时背包 
        BagEquipSlots,      // 装备槽位
        BagItem,            // 物品背包
        MaxBag
    }
} // namespace Phoenix
