using Phoenix.Core;
using System;
using System.Linq;
using System.Collections.Generic;
using UnityEngine;
using System.Text;

namespace Phoenix.Game.FightEmulator
{
    public class ValidList
    {
        public List<eEquipBaseType> valids = new List<eEquipBaseType>();

        public void Add(eEquipBaseType baseType)
        {
            valids.Add(baseType);
        }

        public bool Has(eEquipBaseType baseType)
        {
            for(var i = 0; i < valids.Count; i ++) 
            {
                if (valids[i] == baseType)
                    return true;
            }
            return false;
        }
    }
    public static class EquipUtil
    {
        private static Dictionary<eEquipSlot, ValidList> _slotValids = new Dictionary<eEquipSlot, ValidList>();
        static EquipUtil()
        {
            initSlotValids();
        }

        private static void initSlotValids()
        {
            addValid(eEquipSlot.Head, eEquipBaseType.Head);
            addValid(eEquipSlot.Neck, eEquipBaseType.Neck);
            addValid(eEquipSlot.Back, eEquipBaseType.Back);
            addValid(eEquipSlot.Chest, eEquipBaseType.Chest);
            addValid(eEquipSlot.Wrists, eEquipBaseType.Wrists);
            addValid(eEquipSlot.Hands, eEquipBaseType.Hands);
            addValid(eEquipSlot.Waist, eEquipBaseType.Waist);
            addValid(eEquipSlot.Legs, eEquipBaseType.Legs);
            addValid(eEquipSlot.Feet, eEquipBaseType.Feet);
            addValid(eEquipSlot.Finger1, eEquipBaseType.Finger);
            addValid(eEquipSlot.Finger2, eEquipBaseType.Finger);
            addValid(eEquipSlot.Trinket1, eEquipBaseType.Trinket);
            addValid(eEquipSlot.Trinket2, eEquipBaseType.Trinket);
            addValid(eEquipSlot.MainHand, eEquipBaseType.SingleHand);
            addValid(eEquipSlot.MainHand, eEquipBaseType.DoubleHand);
            addValid(eEquipSlot.MainHand, eEquipBaseType.RangeWeapon);
            addValid(eEquipSlot.OffHand, eEquipBaseType.SingleHand);
            addValid(eEquipSlot.OffHand, eEquipBaseType.Shield);
        }

        private static void addValid(eEquipSlot slot, eEquipBaseType baseType)
        {
            ValidList valids;
            if(!_slotValids.TryGetValue(slot, out valids) )
            {
                valids = new ValidList();
                _slotValids[slot] = valids;
            }

            valids.Add(baseType);
        }

        private static bool hasValid(eEquipSlot slot, eEquipBaseType baseType)
        {
            ValidList valids;
            if (!_slotValids.TryGetValue(slot, out valids))
                return false;
            return valids.Has(baseType);
        }

        public static bool CanEquipInSlot(eEquipSlot slot, eEquipBaseType baseType)
        {            
            if(hasValid(slot, baseType))
                return true;
            return false;
        }

        public static bool HasEquipData(string itemId)
        {
            return EquipDataMgr.It.GetItem(itemId) != null;
        }

        public static bool NeedDoubleHand(eEquipBaseType baseType)
        {
            switch(baseType)
            {
                case eEquipBaseType.DoubleHand:
                case eEquipBaseType.RangeWeapon:
                    return true;
            }
            return false;
        }

        public static EquipItem CreateEquip(string cfgId)
        {
            if (string.IsNullOrEmpty(cfgId))
                return null;
            var equip = new EquipItem();
            equip.Init(cfgId);
            return equip;
        }

        public static void briefAddWeapon(StringBuilder sb, string which, IWeapon w)
        {
            if (w == null)
                return;
            EquipItem weapon = w as EquipItem;
            if (weapon != null && weapon.GetBlock() > 0)
            {
                // 盾牌
                sb.AppendFormat("{0} shield: {1}\n",
                    which, weapon.GetBlock());

            }
            else
            {
                sb.AppendFormat("{0} weapon: ({1}-{2}) speed:{3}\n",
                    which, w.GetMinDmg(), w.GetMaxDmg(), w.GetSpeed());
            }
        }

        public static void tipAddWeapon(StringBuilder sb, IWeapon w)
        {
            if (w == null)
                return;
            EquipItem weapon = w as EquipItem;
            if (weapon != null && weapon.GetBlock() > 0)
            {
                // 盾牌
                sb.AppendFormat("盾牌格挡: {0}\n",
                    weapon.GetBlock());

            }
            else
            {
                sb.AppendFormat("武器伤害: {0}-{1}\n攻速: {2}\n",
                     w.GetMinDmg(), w.GetMaxDmg(), w.GetSpeed());
            }
        }
    }
}// namespace Phoenix
