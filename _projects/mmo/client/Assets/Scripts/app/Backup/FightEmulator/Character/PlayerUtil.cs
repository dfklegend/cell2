using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;
using System.Text;


namespace Phoenix.Game.FightEmulator
{
    public static class PlayerUtil
    {
        public static bool EquipItem(BagSystem.PlayerBags bags, 
            Character playerChar, eEquipSlot slot, EquipItem equip)
        {
            var itemBag = bags.GetBag((int)BagSystem.eBagType.BagItem);
            CheckEquipResult result = new CheckEquipResult();
            if(playerChar.CheckEquip(slot, equip, result))
            {
                // 先从背包移除
                itemBag.RemoveItem(equip.GetIndex());
                return playerChar.Equip(slot, -1, equip);                
            }

            // 根据返回结果
            if (!result.slotOK)
                return false;

            // 检查背包容量
            if(result.needUnequip.Count > 0)
            {
                // 空间不够
                if (itemBag.GetFreeNum() < result.needUnequip.Count - 1)
                    return false;
            }            

            // 先从背包移除
            itemBag.RemoveItem(equip.GetIndex());

            // 卸下需要部位的装备
            for(var i = 0; i < result.needUnequip.Count; i ++)
            {
                var one = result.needUnequip[i];
                var item = playerChar.GetSlotItem(one);
                // unequip it
                playerChar.Unequip(one);
                // move it to bags
                itemBag.AddItem(item);
            }

            var finalSucc = playerChar.Equip(slot, -1, equip);
            if(!finalSucc)
            {
                Log.LogCenter.Default.Error("PlayerUtil.EquipItem finalSucc false!");
            }

            itemBag.DumpInfo();
            return finalSucc;
        }
    }

}// namespace Phoenix
