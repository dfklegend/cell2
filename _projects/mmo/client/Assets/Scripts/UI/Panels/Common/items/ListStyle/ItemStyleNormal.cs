using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Game.FightEmulator.BagSystem;

namespace Phoenix.Game
{    
    // 普通物品
    public class ItemStyleNormal : BaseItemListStyle
    {
        Transform _icon;
        Text _name;
        Text _num;
        BaseItemIconStyle _iconStyle;


        public override void OnReady()
        {
            _icon = TransformUtil.FindComponent<Transform>(_root, "icon");
            _name = TransformUtil.FindComponent<Text>(_root, "name");
            _num = TransformUtil.FindComponent<Text>(_root, "num");
        }

        public override void RefreshInfo(IShowItem v)
        {
            if (v == null)
                return;            

            _name.text = v.GetName();
            _num.text = $"x{v.GetStack()}";

            if (_iconStyle != null)
                _iconStyle.Destroy();
            _iconStyle = ItemIconStyleBuilder.CreateItemStyle(_icon, v, null,
                IconStyleOptions.Tip);
            _iconStyle.RefreshInfo(v, IconStyleOptions.Tip);
        }
    }
} // namespace Phoenix
