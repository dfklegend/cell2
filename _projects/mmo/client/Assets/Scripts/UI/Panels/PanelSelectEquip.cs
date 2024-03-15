using System;
using System.Collections.Generic;
using System.Text;
using Phoenix.Game.FightEmulator;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using static UnityEngine.UI.Dropdown;

namespace Phoenix.Game
{
    [StringType("PanelSelectEquip")]
    public class PanelSelectEquip : BasePanel
    {
        Button _btnOK;
        Button _btnCancel;
        Dropdown _equipType;
        Text _info;
        Text _textTarSlot;

        // 目标slot
        eEquipSlot _tarSlot = eEquipSlot.MainHand;
        string _oldSel = "";
        string _curSel = "";
        Action<string> _cbOK;


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

            _info = TransformUtil.FindComponent<Text>(_root, "BG/labelBrief");
            _textTarSlot = TransformUtil.FindComponent<Text>(_root, "BG/tarSlot");
            _equipType = TransformUtil.FindComponent<Dropdown>(_root, "BG/equipType");            
            _equipType.onValueChanged.AddListener(onEquipTypeSelect);            
        }

        private void BindEvents(bool bind)
        {
            var events = Core.GlobalEvents.It.events;
            events.Bind(EventDefine.InitFight, OnInitFight, bind);
        }

        protected override void onShow()
        {            
            initEquipTypeData();
            trySelectEquip(_oldSel);
        }

        private void onBtnOK()
        {
            _curSel = _equipType.options[_equipType.value].text;
            if (_cbOK != null)
                _cbOK(_curSel);
            Hide();
        }

        private void onBtnCancel()
        {
            Hide();
        }

        private void initEquipTypeData()
        {
            List<OptionData> options = new List<OptionData>();

            var entries = FightEmulator.EquipDataMgr.It.GetAllItem();
            if (entries == null)
                return;
            foreach(var one in entries)
            {
                if (!EquipUtil.CanEquipInSlot(_tarSlot, one.baseType))
                    continue;
                // 检查是否能装备上
                options.Add(new OptionData(one.id));
            }
            _equipType.ClearOptions();
            _equipType.AddOptions(options);
        }

        private void trySelectEquip(string cfgId)
        {
            if (string.IsNullOrEmpty(cfgId))
                return;
            var options = _equipType.options;
            var index = 0;
            for(; index < options.Count; index++)
            {
                var one = options[index];
                if (one.text == cfgId)
                    break;
            }
            if (index >= options.Count)
                    return;
            _equipType.value = index;
        }

        private void onEquipTypeSelect(int sel)
        {
            if (sel < 0 || sel >= _equipType.options.Count)
                return;
            _curSel = _equipType.options[sel].text;

            // refresh infos
            refreshEquipInfo(_curSel);
        }

        private void OnInitFight(params object[] args)
        {   
        }

        public void SelectEquip(eEquipSlot slot, string oldSel, Action<string> cbOK)
        {
            _cbOK = cbOK;
            _tarSlot = slot;
            _oldSel = oldSel;
            refreshTarSlot(slot);
            refreshEquipInfo(_oldSel);
            Show();
        }

        private void refreshTarSlot(eEquipSlot slot)
        {
            if (slot == eEquipSlot.MainHand)
                _textTarSlot.text = "目标: 主手";
            else
                _textTarSlot.text = "目标: 副手";
        }

        private void refreshEquipInfo(string cfgId)
        {
            var equip = EquipUtil.CreateEquip(cfgId);
            StringBuilder sb = new StringBuilder();
            EquipUtil.briefAddWeapon(sb, "", equip);

            _info.text = sb.ToString();
        }
    }
} // namespace Phoenix
