using System.Collections.Generic;
using Phoenix.Utils;

using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game.FightEmulator
{   
    [StringType("PanelCharInfo")]
    public class PanelCharInfo : BasePanel
    {
        Dictionary<eEquipSlot, EquipSlotBox> _slots
            = new Dictionary<eEquipSlot, EquipSlotBox>();
        public override void OnReady()
        {
            SetDepth(PanelDepth.AboveNormal);
            base.OnReady();

            initSlots();
            BindEvents(true);               
        }

        public override void OnDestroy()
        {
            BindEvents(false);
        }        

        private void BindEvents(bool bind)
        {
            var events = Core.GlobalEvents.It.events;            
        }

        protected override void onShow()
        {
            base.onShow();
            refreshEquips();
        }

        private EquipSlotBox createSlot(eEquipSlot slot, string node)
        {
            EquipSlotBox box = new EquipSlotBox();
            box.Init(_root.Find($"BG/{node}"), slot, onSlotClick);
            return box;
        }

        private void addSlot(eEquipSlot slot)
        {
            string node = slot.ToString().ToLower();
            _slots[slot] = createSlot(slot, node);
        }

        private void initSlots()
        {
            addSlot(eEquipSlot.MainHand);
            addSlot(eEquipSlot.OffHand);
        }

        private EquipSlotBox getBox(eEquipSlot slot)
        {
            EquipSlotBox box;
            if (_slots.TryGetValue(slot, out box))
                return box;
            return null;
        }

        private void refreshEquips()
        {
            var player = FightCtrl.It.playerChar;
            if (player == null)
                return;
            var slots = player.GetSlots();
            for(var i = 0; i < slots.Length; i ++)
            {
                var one = slots[i];                
                var box = getBox(one.index);
                if (box == null)
                    continue;
                box.RefreshInfo(one.item);
            }
        }

        private void onSlotClick(EquipSlotBox box)
        {
            Log.LogCenter.Default.Debug($"box {box.slot.ToString()}");            
            
            if(box.item == null)
            {
                selectEquipToReplace(box);
                return;
            }
            //var panel = TipUtil.ShowEquip(eEquipTipShowType.EquipSlot,
            //    box.slot, box.item,
            //    ()=> {
            //        Debug.Log("select replace");
            //        selectEquipToReplace(box);
            //    });            
        }

        private void selectEquipToReplace(EquipSlotBox box)
        {
            //BagCtrl.It.SelectEquip(box.slot, box.item, (result)=> 
            //{
            //    Debug.Log("SelectEquip ret:" + result.ToString());
            //    if(result.selected != null)
            //    {
            //        EquipItem equip = result.selected as EquipItem;
            //        if (equip == null)
            //            return;
            //        var playerChar = FightCtrl.It.playerChar;

            //        var succ = PlayerUtil.EquipItem(FightCtrl.It.player.bags,
            //            playerChar, box.slot, equip);                    

            //        BagCtrl.It.HideSelectPanel();
            //        refreshEquips();
            //        return;
            //    }

            //    if(result.selectUnequip)
            //    {
            //        var playerChar = FightCtrl.It.playerChar;
            //        playerChar.Unequip(box.slot);
                    
            //        BagCtrl.It.HideSelectPanel();
            //        refreshEquips();
            //    }
            //});
        }

        private void tryEquip(eEquipSlot slot, EquipItem equip)
        {
            // 首先检查是否允许装备
        }
    }
} // namespace Phoenix
