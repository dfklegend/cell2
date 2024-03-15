using System.Collections.Generic;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;

namespace Phoenix.Game
{    
    public static class ItemBagUtil
    {
        // 根据显示物品构建直接的显示列表
        public static  void BuildItems(List<BaseItemListStyle> items, BagEnvData data, Transform root)
        {
            if (data == null)
                return;
            var viewData = data.viewData;            
            var handler = data.handler;
            for (var i = 0; i < viewData.GetSize(); i++)
            {
                var one = viewData.Get(i);
                if (one.ViewType == eBagViewItemType.Label)
                {
                    var label = one as ViewItemLabel;
                    var style = ItemListStyleBuilder.CreateLabel(root, label.label, label.op, handler);
                    items.Add(style);
                    continue;
                }

                if (one.ViewType == eBagViewItemType.Item)
                {
                    var item = one as ViewItem;
                    var style = ItemListStyleBuilder.CreateItemStyle(root, item.item, handler);
                    items.Add(style);
                    continue;
                }
            }
        }
    }
} // namespace Phoenix
