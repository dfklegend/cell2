using System.Collections.Generic;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Game.FightEmulator;
using Phoenix.Game.FightEmulator.BagSystem;
using Phoenix.Core;
using System;

namespace Phoenix.Game
{   
    public class SelectEquipResult
    {   
        public bool selectUnequip;        
        public IShowItem selected;
    }

    // 背包显示操作
    public class BagCtrl : Singleton<BagCtrl>
    {
        // 显示bag
        public void ShowItemBag()
        {       
            var envData = new BagEnvData();
            envData.viewData = makeItemData();// makeTestData("test");
            envData.handler = new EquipBagHandler();

            UIMgr.It.GetPanel<PanelItemBag>("PanelItemBag").ShowBag(envData);
        }


        private BagViewData makeTestData(string title)
        {
            BagViewData data = new BagViewData();

            var item = FightEmulator.ItemSystem.ItemBuilder.CreateItem("短剑_1");
            data.Add(new ViewItemLabel($"{title}开始", eBagLabelOp.Empty));
            //data.Add(new ViewItem(item));
            //data.Add(new ViewItem(item));
            data.Add(new ViewItemLabel($"{title}中间", eBagLabelOp.Empty));
            //for(var i = 0; i < 20; i ++)
            //    data.Add(new ViewItem(item));
            data.Add(new ViewItemLabel($"{title}结束", eBagLabelOp.Empty));
            return data;
        }

        private BagViewData makeItemData()
        {
            BagViewData data = new BagViewData();            
            return data;
        }

        public void HideSelectPanel()
        {
            UIMgr.It.GetPanel<PanelBagSelect>("PanelBagSelect").Hide();
        }

        public void ShowEquipBag()
        {
            var envData = new BagEnvData();
            envData.viewData = makeEquipToSelect(0);
            envData.handler = new EquipBagHandler();

            UIMgr.It.GetPanel<PanelItemBag>("PanelItemBag").ShowBag(envData);
        }


        public void SelectEquip(int index, Action<SelectEquipResult> cb)
        {
            var envData = new BagEnvData();
            envData.viewData = makeEquipToSelect(index);

            var handler = new SelectEquipHandler(cb);
            envData.handler = handler;

            UIMgr.It.GetPanel<PanelBagSelect>("PanelBagSelect").ShowBag("选择装备", envData);
        }

        private BagViewData makeEquipToSelect(int index)
        {
            BagViewData data = new BagViewData();

            // 每一个道具
            var equips = Card.EquipDataMgr.It.GetAllItem();
            foreach(var one in equips)
            {
                if (string.IsNullOrEmpty(one.id))
                    continue;
                data.Add(new ViewItem(new Card.EquipItem(one.id)));
            }

            return data;
        }

        public void SelectEquipToReplace(int index, IShowItem item, Action<SelectEquipResult> cb)
        {
            var envData = new BagEnvData();
            envData.viewData = makeEquipToReplace(index, item);

            var handler = new SelectEquipHandler(cb);            
            envData.handler = handler;

            UIMgr.It.GetPanel<PanelBagSelect>("PanelBagSelect").ShowBag("替换装备", envData);
        }        

        private BagViewData makeEquipToReplace(int index, IShowItem item)
        {
            BagViewData data = new BagViewData();      

            data.Add(new ViewItemLabel($"卸下", eBagLabelOp.Unequip));
            var equips = Card.EquipDataMgr.It.GetAllItem();
            foreach (var one in equips)
            {
                if (string.IsNullOrEmpty(one.id))
                    continue;
                data.Add(new ViewItem(new Card.EquipItem(one.id)));
            }

            return data;
        }

        public void ShowSkillBag()
        {
            var envData = new BagEnvData();
            envData.viewData = makeSkillToSelect(0);
            envData.handler = new SkillBagHandler();

            UIMgr.It.GetPanel<PanelItemBag>("PanelItemBag").ShowBag(envData);
        }

        public void SelectSkill(int index, Action<SelectEquipResult> cb)
        {
            var envData = new BagEnvData();
            envData.viewData = makeSkillToSelect(index);

            var handler = new SelectEquipHandler(cb);
            envData.handler = handler;

            UIMgr.It.GetPanel<PanelBagSelect>("PanelBagSelect").ShowBag("选择技能", envData);
        }

        public void SelectSkillToReplace(int index, IShowItem item, Action<SelectEquipResult> cb)
        {
            var envData = new BagEnvData();
            envData.viewData = makeSkillToReplace(index, item);

            var handler = new SelectEquipHandler(cb);
            envData.handler = handler;

            UIMgr.It.GetPanel<PanelBagSelect>("PanelBagSelect").ShowBag("替换技能", envData);
        }

        private bool isSkillValid(Card.SkillData one)
        {
            if (string.IsNullOrEmpty(one.id))
                return false;
            // 去掉普攻
            if (one.normalAttack > 0)
                return false;
            if (one.canSelect == 0)
                return false;
            return true;
        }

        private BagViewData makeSkillToSelect(int index)
        {
            BagViewData data = new BagViewData();

            // 每一个道具
            var skills = Card.SkillDataMgr.It.GetAllItem();
            foreach (var one in skills)
            {
                if (!isSkillValid(one))
                    continue;                
                data.Add(new ViewItem(new Card.SkillItem(one.id)));
            }

            return data;
        }

        private BagViewData makeSkillToReplace(int index, IShowItem item)
        {
            BagViewData data = new BagViewData();

            data.Add(new ViewItemLabel($"卸下", eBagLabelOp.Unequip));
            // 每一个道具
            var skills = Card.SkillDataMgr.It.GetAllItem();
            foreach (var one in skills)
            {
                if (!isSkillValid(one))
                    continue;
                if (item != null && one.id == item.GetItemId())
                    continue;
                data.Add(new ViewItem(new Card.SkillItem(one.id)));
            }

            return data;
        }
    }
} // namespace Phoenix
