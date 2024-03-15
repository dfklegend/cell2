using Phoenix.Core;
using System;
using System.Linq;
using System.Collections.Generic;
using UnityEngine;
using System.Text;

namespace Phoenix.Game.FightEmulator
{
    public class EquipSlot
    {
        public eEquipSlot index;
        // 装备的物品
        public EquipItem item;
    }

    public class CheckEquipResult
    {
        public bool allOK = false;
        public bool slotOK = false;
        public List<eEquipSlot> needUnequip = new List<eEquipSlot>();

        public void Reset()
        {
            allOK = false;
            slotOK = false;
            needUnequip.Clear();
        }
    }

    // 不应该关心背包
    // 装备前，可以提供检查，获取
    //     是否部位合法
    //     如果装备需要卸下的物品
    // 由外部保证后，才允许装备
    public class EquipComponent
    {
        Character _character;
        private EquipSlot[] _slots = new EquipSlot[(int)eEquipSlot.Max];
        public EquipSlot[] slots { get { return _slots; } }

        public EquipComponent()
        {
            for (var i = 0; i < (int)eEquipSlot.Max; i++)
            {
                _slots[i] = new EquipSlot();
                _slots[i].index = (eEquipSlot)i;
            }
                
        }

        public void Init(Character c)
        {
            _character = c;
        }

        private EquipSlot getSlot(eEquipSlot slotIndex)
        {
            if (slotIndex < 0 || slotIndex >= eEquipSlot.Max)
                return null;
            return _slots[(int)slotIndex];
        }

        public bool EquipItem(eEquipSlot slotIndex, int indexInBag, EquipItem item)
        {
            if (item == null)
                return false;
            var slot = getSlot(slotIndex);
            if(slot == null)
            {
                Log.LogCenter.Default.Warning("error slot:", slotIndex);
                return false;
            }

            CheckEquipResult result = new CheckEquipResult();
            CheckEquipItem(slotIndex, item, result);
            if(!result.allOK)
            {
                Log.LogCenter.Default.Warning("can now equip at {0}", slotIndex);
                return false;
            }

            slot.item = item;
            _character.ApplyEquip(slotIndex, item);            
            return true;
        }        

        public void CheckEquipItem(eEquipSlot slotIndex, EquipItem item, CheckEquipResult result)
        {
            result.allOK = false;
            if (item == null)
                return;
            var slot = getSlot(slotIndex);
            if (slot == null)
                return;

            // check whether this can equip in slot
            result.slotOK = EquipUtil.CanEquipInSlot(slotIndex, item.GetEquipBaseType());

            EquipItem oldItem = slot.item;
            if (oldItem != null)
            {
                result.needUnequip.Add(slotIndex);
            }

            // 特殊: 双手武器副手也要卸下
            if (slotIndex == eEquipSlot.MainHand)
            {
                bool doubleHand = EquipUtil.NeedDoubleHand(item.GetEquipBaseType());
                if (doubleHand)
                {
                    tryAddNeedUnequip(result, eEquipSlot.OffHand);
                }
            }

            // 装备副手，主手如果是双手也卸下
            if (slotIndex == eEquipSlot.OffHand)
            {
                var mainHand = getSlot(eEquipSlot.MainHand);
                if (mainHand.item != null && EquipUtil.NeedDoubleHand(mainHand.item.GetEquipBaseType()))
                    result.needUnequip.Add(eEquipSlot.MainHand);
            }

            result.allOK = result.slotOK && result.needUnequip.Count == 0;
        }

        private void tryAddNeedUnequip(CheckEquipResult result, eEquipSlot slotIndex)
        {
            if (getSlot(slotIndex).item == null)
                return;
            result.needUnequip.Add(slotIndex);
        }


        public bool UnequipItem(eEquipSlot slotIndex)
        {
            var slot = getSlot(slotIndex);
            if (slot == null)
            {
                Log.LogCenter.Default.Warning("error slot:", slotIndex);
                return false;
            }
            
            if (slot.item != null)
            {
                unequipSlot(slot, -1);
                return true;
            }
            return false;
        }        

        private void unequipSlot(EquipSlot slot, int indexInBag)
        {
            if (slot.item == null)
                return;
            
            var oldItem = slot.item;
            slot.item = null;
            _character.ApplyUnequip(slot.index, oldItem);            
        }

        public EquipItem GetSlotItem(eEquipSlot slotIndex)
        {
            var slot = getSlot(slotIndex);
            if (slot != null)
                return slot.item;
            return null;
        }

        public void Dump()
        {   
            StringBuilder sb = new StringBuilder();
            MakeBrief(sb);
            Log.LogCenter.Default.Debug(sb.ToString());
        }

        public void MakeBrief(StringBuilder sb)
        {            
            for (var i = 0; i < _slots.Length; i++)
            {
                var one = _slots[i];
                if (one.item != null)
                {
                    sb.AppendLine($"slot:{(eEquipSlot)i} item:{one.item.GetCfgId()}");
                    //Log.LogCenter.Default.Debug($"slot:{(eEquipSlot)i} item:{one.item.GetCfgId()}");
                }
            }
        }
    }
}// namespace Phoenix
