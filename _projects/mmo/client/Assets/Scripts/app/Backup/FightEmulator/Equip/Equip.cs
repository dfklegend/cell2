using UnityEngine;
using System.Reflection;
using Phoenix.csv;
using Phoenix.Game.FightEmulator.ItemSystem;

namespace Phoenix.Game.FightEmulator
{
    public class EquipItem : BaseItem, IWeapon
    {
        private ItemData _itemData;
        private EquipData _equipData;
        // 后续
        // 随机属性


        public override void Init(string cfgId)
        {
            base.Init(cfgId);
            _itemData = _cfg;
            _equipData = EquipDataMgr.It.GetItem(cfgId);
        }

        public string GetCfgId()
        {
            return _cfgId;
        }

        public eEquipBaseType GetEquipBaseType()
        {
            return (eEquipBaseType)_equipData.baseType;
        }
        
        public int GetMinDmg()
        {
            return _equipData.minDmg;
        }

        public int GetMaxDmg()
        {
            return _equipData.maxDmg;
        }

        public float GetSpeed()
        {
            return _equipData.speed;
        }

        public int GetBlock()
        {
            return _equipData.block;
        }

        public void Reset()
        { }

        public void Equip(IAttrOwner owner)
        {
            applyEquip(owner, 1f);
        }

        public void Unequip(IAttrOwner owner)
        {
            applyEquip(owner, -1f);
        }

        private void applyEquip(IAttrOwner owner, float factor)
        {
            if (_equipData == null)
                return;            
            var d = _equipData;
            if (d.value0 > 0)
                tryAddAttr(owner, d.attr0, d.value0 * factor);
            if (d.value1 > 0)
                tryAddAttr(owner, d.attr1, d.value1 * factor);
            if (d.value2 > 0)
                tryAddAttr(owner, d.attr2, d.value2 * factor);
            if (d.value3 > 0)
                tryAddAttr(owner, d.attr3, d.value3 * factor);

            owner.RecalcAttrTransform();
        }

        private void tryAddAttr(IAttrOwner owner, string typeString, float value)
        {            
            AttrsUtil.AddEquipAttr(owner, typeString, value);
        }

        public string MakeBrief()
        {
            return "";
        }
    }
}
