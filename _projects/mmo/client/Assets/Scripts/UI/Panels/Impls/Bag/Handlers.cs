using System.Collections.Generic;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Game.FightEmulator.ItemSystem;
using Phoenix.Game.FightEmulator.BagSystem;
using System;
using Phoenix.Game.FightEmulator;

namespace Phoenix.Game
{
    // 背包内，点击显示tip
    public class EquipBagHandler : IShowItemHandler
    {
        public void OnClick(IShowItemStyle style, IShowItem item)
        {
            if (item == null)
                return;
            // 显示物品tip
            Debug.Log("Show tip:" + item.GetItemId());            
            TipUtil.ShowEquip(eEquipTipShowType.Bag,item);
        }
    }


    public class SelectEquipHandler : IShowItemHandler
    {
        Action<SelectEquipResult> _cb;

        public SelectEquipHandler(Action<SelectEquipResult> cb)
        {
            _cb = cb;
        }

        public void OnClick(IShowItemStyle style, IShowItem item)
        {
            SelectEquipResult result = new SelectEquipResult();
            // 卸下
            if(item == null)
            {
                ItemStyleLabel label = style as ItemStyleLabel;
                if(label != null)
                {
                    result.selectUnequip = label.op == eBagLabelOp.Unequip;
                }
            }
            // 装备
            result.selected = item;

            Debug.Log(result);
            _cb?.Invoke(result);
        }
    }

    public class SkillBagHandler : IShowItemHandler
    {
        public void OnClick(IShowItemStyle style, IShowItem item)
        {
            if (item == null)
                return;
            // 显示物品tip
            Debug.Log("Show tip:" + item.GetItemId());
            TipUtil.ShowSkill(eSkillTipShowType.Bag, item);
        }
    }
} // namespace Phoenix
