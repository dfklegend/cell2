using Phoenix.Game.FightEmulator.BagSystem;
using Phoenix.Utils;
using System;
using System.Text;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game
{
    // tip界面，底下功能按钮根据情况显示不同内容

    public enum eEquipTipShowType
    {
        Equip,
        Bag
    }

    public class EquipTipData
    {
        public int slot;
        public eEquipTipShowType showType;
        public IShowItem item;
    }    

    public class FuncBtn
    {
        public Button btn;
        public Text text;
    }

    [StringType("PanelEquipTip")]
    public class PanelEquipTip : BasePanel
    {
        const int MAX_BTN = 3;
        Transform _icon;
        Text _name;
        Text _info;
        BaseItemIconStyle _style;

        EquipTipData _data;

        FuncBtn[] _btns = new FuncBtn[MAX_BTN];
        Action[] _handlers = new Action[MAX_BTN];        

        public override void OnReady()
        {
            SetDepth(PanelDepth.AboveNormal + 2);
            base.OnReady();

            _icon = TransformUtil.FindComponent<Transform>(_root, "BG/base/icon");
            _name = TransformUtil.FindComponent<Text>(_root, "BG/base/name");
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

        public void ShowTip(EquipTipData data)
        {
            _data = data;
            showAttrs();
            Show();
        }

        private void showAttrs()
        {
            if (_data == null || _data.item == null)
                return;
            var cfg = Card.EquipDataMgr.It.GetItem(_data.item.GetItemId());
            if (cfg == null)
                return;
            var sb = new StringBuilder();
            sb.Append("属性:\n");            
            makeAttr(sb, cfg.attr0, cfg.v0);
            makeAttr(sb, cfg.attr1, cfg.v1);
            makeAttr(sb, cfg.attr2, cfg.v2);
            makeAttr(sb, cfg.attr3, cfg.v3);
            _info.text = sb.ToString();
        }

        private void makeAttr(StringBuilder sb, Card.AttrType at, float v )
        {
            if (at.type != eAttrType.Invalid)
                Card.AttrFormatter.FormatAttr(sb, "    ", at.type, at.percent, v);
        }

        private void refreshInfo(IShowItem item)
        {
            if (item == null)
                return;
            _style.RefreshInfo(item, IconStyleOptions.Tip);
            _name.text = item.GetName();

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

        public void ShowBtn(int num)
        {
            for(var i = 0; i < MAX_BTN; i ++)
            {
                bool visible = i < num;
                _btns[i].btn.gameObject.SetActive(visible);    
            }
        }

        public void ResetHandlers()
        {
            for (var i = 0; i < _handlers.Length; i++)
                _handlers[i] = null;
        }

        public void SetBtnText(int index, string text)
        {
            _btns[index].text.text = text;
        }
    }
} // namespace Phoenix
