using Phoenix.Game.FightEmulator.BagSystem;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game.FightEmulator
{
    public class EquipSlotBox
    {
        public delegate void OnClickSlot(EquipSlotBox box);

        // 具体的装备部位
        eEquipSlot _slot;
        public eEquipSlot slot { get { return _slot; } }
        EquipItem _item;
        public EquipItem item { get { return _item; } }

        Transform _root;
        Image _icon;
        Text _name;
        Button _btn;
        OnClickSlot _handler;

        BaseItemIconStyle _style;

        public void Init(Transform root,
            eEquipSlot slot, OnClickSlot clickHandler)
        {
            _slot = slot;
            if (!root)
                return;
            _root = root;
            _handler = clickHandler;



            //_icon = TransformUtil.FindComponent<Image>(_root, "icon");
            //_name = TransformUtil.FindComponent<Text>(_root, "name");
            //_btn = TransformUtil.FindComponent<Button>(_root, "btn");

            //_btn.onClick.AddListener(onClick);
            
        }

        public void RefreshInfo(IBagItem v)
        {
            //_item = v as EquipItem;
            //if (_style == null)
            //    _style = IconStyleBuilder.CreateItemStyle(eItemStyle.Normal, _root, v, onClickBtn, null);
            //else
            //    _style.RefreshInfo(v, null);
        }

        private void resetInfo()
        {
            _icon.sprite = null;
            _name.text = "";
        }

        private void onClick()
        {
            if (_handler != null)
                _handler(this);
        }

        private void onClickBtn(IBagItem v, object arg)
        {
            if (_handler != null)
                _handler(this);
        }
    }    
} // namespace Phoenix
