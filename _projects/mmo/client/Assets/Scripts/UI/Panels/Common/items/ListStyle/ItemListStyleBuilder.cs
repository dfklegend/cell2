using System.Collections.Generic;
using Phoenix.Core;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Game.FightEmulator.BagSystem;

namespace Phoenix.Game
{    
    // 物品的list样式
    public static class ItemListStyleBuilder
    {
        public static BaseItemListStyle CreateItemStyle(Transform node, IShowItem item, IShowItemHandler handler)
        {
            var styleType = GetItemStyle(item);
            var style = ItemListStyleFactory.It.Create(styleType);
            if (style == null)
                return null;
            var go = ItemPrefabFactory.It.ListInstantiate(styleType);
            if (go == null)
                return null;
            go.transform.SetParent(node, false);
            
            var holder = go.GetComponent<ItemStyleHolder>();
            
            holder.style = style;
            holder.item = item;            

            style.Init(go);
            style.RefreshInfo(item);
            style.handler = handler;
            return style;
        }

        public static BaseItemListStyle CreateLabel(Transform node, string text, 
            eBagLabelOp op, IShowItemHandler handler)
        {
            var styleType = eItemStyle.Label;
            var style = ItemListStyleFactory.It.Create(styleType) as ItemStyleLabel;
            if (style == null)
                return null;
            var go = ItemPrefabFactory.It.ListInstantiate(styleType);
            if (go == null)
                return null;
            go.transform.SetParent(node, false);

            var holder = go.GetComponent<ItemStyleHolder>();

            holder.style = style;
            holder.item = null;

            style.Init(go);
            style.SetInfo(text);

            style.handler = handler;
            style.op = op;
            return style;
        }

        public static eItemStyle GetItemStyle(IShowItem item)
        {
            return eItemStyle.Normal;
        }
    }
    
    
} // namespace Phoenix
