using Phoenix.Game.FightEmulator.BagSystem;
using Phoenix.Utils;
using System;
using System.Text;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game
{
    public enum eSkillTipShowType
    {
        Equip,
        Bag
    }

    public class SkillTipData
    {
        public int slot;
        public eSkillTipShowType showType;
        public IShowItem item;
    }    
   

    [StringType("PanelSkillTip")]
    public class PanelSkillTip : BasePanel
    {
        const int MAX_BTN = 3;
        Transform _icon;
        Text _name;
        Text _cd;
        Text _info;
        BaseItemIconStyle _style;

        SkillTipData _data;

        Card.SkillData _cfg;

        FuncBtn[] _btns = new FuncBtn[MAX_BTN];
        Action[] _handlers = new Action[MAX_BTN];        

        public override void OnReady()
        {
            SetDepth(PanelDepth.AboveNormal + 2);
            base.OnReady();

            _icon = TransformUtil.FindComponent<Transform>(_root, "BG/base/icon");
            _name = TransformUtil.FindComponent<Text>(_root, "BG/base/name");
            _cd = TransformUtil.FindComponent<Text>(_root, "BG/base/cd");
            _info = TransformUtil.FindComponent<Text>(_root, "BG/info/text");

            _style = ItemIconStyleBuilder.CreateItemStyle(eItemStyle.Normal, _icon, null, null,
                IconStyleOptions.Tip);
            initBtns();
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
            if (_data == null)
                return;
            refreshInfo(_data.item);
        }

        public void ShowTip(SkillTipData data)
        {
            _data = data;
            Show();
        }

        private void refreshInfo(IShowItem item)
        {
            if (item == null)
                return;
            _cfg = Card.SkillDataMgr.It.GetItem(item.GetItemId());
            if (_cfg == null)
                return;
            _style.RefreshInfo(item, IconStyleOptions.Tip);
            _name.text = item.GetName();
            _cd.text = ""+_cfg.cd;
            _info.text = _cfg.desc;

            //StringBuilder sb = new StringBuilder();
            //EquipUtil.tipAddWeapon(sb, item);

            //_info.text = sb.ToString();
        }

        private void initBtns()
        {
            UnityEngine.Events.UnityAction[] actions = { onBtn0, onBtn1, onBtn2 };

            for(var i = 0; i < MAX_BTN; i ++)
            {
                initBtn(i, actions[i]);
            }            
        }

        private void initBtn(int index, UnityEngine.Events.UnityAction cb)
        {
            var btn = TransformUtil.FindComponent<Button>(_root, $"funcs/btn{index}");

            btn.onClick.AddListener(cb);

            var funcBtn = new FuncBtn();
            _btns[index] = funcBtn;
            funcBtn.btn = btn;
            funcBtn.text = TransformUtil.FindComponent<Text>(btn.transform, "Text");
        }

        private void onBtn0()
        {
            Hide();
            tryInvoke(0);
        }

        private void onBtn1()
        {
            Hide();
            tryInvoke(1);
        }

        private void onBtn2()
        {
            Hide();
            tryInvoke(2);
        }

        public void SetHandler(int index, Action cb)
        {
            if (cb == null)
                return;
            if (index < 0 || index >= 3)
                return;
            _handlers[index] = cb;
        }

        private void tryInvoke(int index)
        {
            var action = _handlers[index];
            if (action == null)
                return;
            action();
        }

        public void ResetHandlers()
        {
            for (var i = 0; i < _handlers.Length; i++)
                _handlers[i] = null;
        }

        public void ShowBtn(int num)
        {
            for(var i = 0; i < MAX_BTN; i ++)
            {
                bool visible = i < num;
                _btns[i].btn.gameObject.SetActive(visible);    
            }
        }

        public void SetBtnText(int index, string text)
        {
            _btns[index].text.text = text;
        }
    }
} // namespace Phoenix
