using System.Collections.Generic;
using Phoenix.Game.FightEmulator;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using static UnityEngine.UI.Dropdown;

namespace Phoenix.Game
{
    [StringType("PanelSelectEnemy")]
    public class PanelSelectEnemy : BasePanel
    {
        Button _btnOK;
        Button _btnCancel;
        Dropdown _enemyType;
        InputField _inputLevel;
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
            _enemyType = TransformUtil.FindComponent<Dropdown>(_root, "BG/enemyType");            
            _enemyType.onValueChanged.AddListener(onEnemyTypeSelect);
            _inputLevel = TransformUtil.FindComponent<InputField>(_root, "BG/inputLevel");
        }

        private void BindEvents(bool bind)
        {
            var events = Core.GlobalEvents.It.events;
            events.Bind(EventDefine.InitFight, OnInitFight, bind);
        }

        protected override void onShow()
        {
            _inputLevel.text = "" + FightSimulateCtrl.It.enemyLevel;
        }

        private void onBtnOK()
        {
            FightSimulateCtrl.It.enemyCfgId = _enemyType.options[_enemyType.value].text;
            FightSimulateCtrl.It.enemyLevel = int.Parse(_inputLevel.text);
            Hide();
        }

        private void onBtnCancel()
        {
            Hide();
        }

        private void initEnemyTypeData()
        {
            List<OptionData> options = new List<OptionData>();

            var enemies = FightEmulator.MonsterDataMgr.It.GetAllItem();
            if (enemies == null)
                return;
            foreach(var one in enemies)
            {
                options.Add(new OptionData(one.id));
            }
            _enemyType.ClearOptions();
            _enemyType.AddOptions(options);
        }

        private void onEnemyTypeSelect(int sel)
        {
        }

        private void OnInitFight(params object[] args)
        {
            initEnemyTypeData();
        }

        public void ReqSelect(string cfgId)
        {
            trySelect(cfgId);
            Show();
        }

        private void trySelect(string cfgId)
        {
            if (string.IsNullOrEmpty(cfgId))
                return;
            var options = _enemyType.options;
            var index = 0;
            for (; index < options.Count; index++)
            {
                var one = options[index];
                if (one.text == cfgId)
                    break;
            }
            if (index >= options.Count)
                return;
            _enemyType.value = index;
        }
    }
} // namespace Phoenix
