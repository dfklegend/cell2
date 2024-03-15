using System.Collections.Generic;
using Phoenix.Game.FightEmulator;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.Events;
using UnityEngine.UI;
using static UnityEngine.UI.Dropdown;

namespace Phoenix.Game
{
    public class Slot
    {
        public Text itemName;
        public Button btn;
    }

    public enum eSlot
    {
        MainHand = 0,
        OffHand,
        Max
    }

    [StringType("PanelSelectRole")]
    public class PanelSelectRole : BasePanel
    {
        Button _btnOK;
        Button _btnCancel;      
        InputField _inputLevel;

        List<Slot> _slots = new List<Slot>();

        public override void OnReady()
        {            
            base.OnReady();

            initCtrls();
            BindEvents(true);
        }

        public override void OnDestroy()
        {
            BindEvents(false);
        }

        private void initCtrls()
        {
            _btnOK = TransformUtil.FindComponent<Button>(_root, "BG/btnOK");
            _btnOK.onClick.AddListener(onBtnOK);
            _btnCancel = TransformUtil.FindComponent<Button>(_root, "BG/btnCancel");
            _btnCancel.onClick.AddListener(onBtnCancel);            
            _inputLevel = TransformUtil.FindComponent<InputField>(_root, "BG/inputLevel");

            addSlot(_root.Find("BG/mainhand"), OnMainHand);
            addSlot(_root.Find("BG/offhand"), OnOffHand);
        }

        private void addSlot(Transform node, UnityAction action)
        {
            _slots.Add(initSlot(node, action));
        }

        private Slot initSlot(Transform node, UnityAction action)
        {
            var slot = new Slot();
            slot.itemName = TransformUtil.FindComponent<Text>(node, "btn/Text");
            slot.btn = TransformUtil.FindComponent<Button>(node, "btn");
            slot.btn.onClick.AddListener(action);
            return slot;
        }

        private void BindEvents(bool bind)
        {
            var events = Core.GlobalEvents.It.events;
            events.Bind(EventDefine.InitFight, OnInitFight, bind);
        }

        protected override void onShow()
        {
            _inputLevel.text = "" + FightSimulateCtrl.It.roleLevel;
        }

        private void onBtnOK()
        {
            FightSimulateCtrl.It.roleLevel = int.Parse(_inputLevel.text);
            Hide();
        }

        private void onBtnCancel()
        {
            Hide();
        }

        private void OnInitFight(params object[] args)
        {
            getSlot(eSlot.MainHand).itemName.text = FightSimulateCtrl.It.mainHand;
            getSlot(eSlot.OffHand).itemName.text = FightSimulateCtrl.It.offHand;
            _inputLevel.text = "" + FightSimulateCtrl.It.roleLevel;
        }        

        private Slot getSlot(eSlot index)
        {
            if (index < 0 || index >= eSlot.Max)
                return null;
            return _slots[(int)index];
        }

        private void OnMainHand()
        {
            Debug.Log("OnMainHand");
            UIMgr.It.GetPanel<PanelSelectEquip>("PanelSelectEquip").SelectEquip(eEquipSlot.MainHand,
                FightSimulateCtrl.It.mainHand, (newSel) => 
                {
                    FightSimulateCtrl.It.mainHand = newSel;
                    makeSureEquipValid(true);
                });
        }

        private void OnOffHand()
        {
            Debug.Log("OnOffHand");
            UIMgr.It.GetPanel<PanelSelectEquip>("PanelSelectEquip").SelectEquip(eEquipSlot.OffHand,
                FightSimulateCtrl.It.offHand, (newSel) =>
                {
                    FightSimulateCtrl.It.offHand = newSel;
                    makeSureEquipValid(false);
                });
        }

        private void makeSureEquipValid(bool equipMain)
        {
            var main = FightSimulateCtrl.It.mainHand;
            var off = FightSimulateCtrl.It.offHand;

            var c = new Character();
            var mainEquip = EquipUtil.CreateEquip(main);
            var offEquip = EquipUtil.CreateEquip(off);

            if(equipMain)
            {
                c.Equip(eEquipSlot.OffHand, -1, offEquip);
                c.Equip(eEquipSlot.MainHand, -1, mainEquip);                
            }
            else
            {
                c.Equip(eEquipSlot.MainHand, -1, mainEquip);
                c.Equip(eEquipSlot.OffHand, -1, offEquip);
            }
            

            var finalMain = c.GetMainWeapon() as EquipItem;
            if (finalMain != null)
            {
                FightSimulateCtrl.It.mainHand = finalMain.GetCfgId();
            }
            else
                FightSimulateCtrl.It.mainHand = "";
            var finalOff = c.GetOffHandWeapon() as EquipItem;
            if (finalOff != null)
                FightSimulateCtrl.It.offHand = finalOff.GetCfgId();
            else
                FightSimulateCtrl.It.offHand = "";

            getSlot(eSlot.MainHand).itemName.text = FightSimulateCtrl.It.mainHand;
            getSlot(eSlot.OffHand).itemName.text = FightSimulateCtrl.It.offHand;
        }
    }
} // namespace Phoenix
