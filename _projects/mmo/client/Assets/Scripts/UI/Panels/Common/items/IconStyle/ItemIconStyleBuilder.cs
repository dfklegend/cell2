using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Game.FightEmulator.BagSystem;
using System;

namespace Phoenix.Game
{
    // 物品的icon样式
    public static class ItemIconStyleBuilder
    {
        public static BaseItemIconStyle CreateItemStyle(Transform node, IShowItem item,
            Action<IShowItem, object> handler, IIconStyleOptions options)
        {
            var styleType = GetItemStyle(item);
            return CreateItemStyle(styleType, node, item, handler, options);
        }

        public static BaseItemIconStyle CreateItemStyle(eItemStyle styleType, 
            Transform node, IShowItem item, Action<IShowItem, object> handler,
            IIconStyleOptions options)
        {   
            var style = ItemIconStyleFactory.It.Create(styleType);
            if (style == null)
                return null;
            var go = ItemPrefabFactory.It.IconInstantiate(styleType);
            if (go == null)
                return null;
            go.transform.SetParent(node, false);

            style.Init(go);
            style.RefreshInfo(item, options);
            style.SetHandler(handler, null);
            return style;
        }

        public static eItemStyle GetItemStyle(IShowItem item)
        {
            return eItemStyle.Normal;
        }
    }
} // namespace Phoenix
