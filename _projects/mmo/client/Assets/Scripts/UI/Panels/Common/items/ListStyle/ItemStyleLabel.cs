using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Game.FightEmulator.BagSystem;

namespace Phoenix.Game
{   
    public enum eBagLabelOp
    {
        Empty, 
        Unequip,            // 卸下
    }

    // 抽象物品显示
    public class ItemStyleLabel : BaseItemListStyle
    {
        Text _text;
        // 存储用于后面
        public eBagLabelOp op = eBagLabelOp.Empty;

        public override void OnReady()
        {
            _text = TransformUtil.FindComponent<Text>(_root, "text");            
        }

        public void SetInfo(string v)
        {
            _text.text = v;
        }

        public override void RefreshInfo(IShowItem v)
        {            
        }
    }
} // namespace Phoenix
