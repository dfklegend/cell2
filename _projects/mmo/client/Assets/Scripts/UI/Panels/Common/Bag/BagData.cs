using System.Collections.Generic;
using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Game.FightEmulator.ItemSystem;
using Phoenix.Game.FightEmulator.BagSystem;
using System;

namespace Phoenix.Game
{
    public enum eBagViewItemType
    {
        Label,
        Item
    }

    public interface IBagViewItem
    {
        eBagViewItemType ViewType { get; }
    }

    public class ViewItemLabel : IBagViewItem
    {
        public eBagViewItemType ViewType { get { return eBagViewItemType.Label; } }

        public string label;
        public eBagLabelOp op;

        public ViewItemLabel(string text, eBagLabelOp op)
        {
            label = text;
            this.op = op;
        }
    }

    public class ViewItem : IBagViewItem
    {
        public eBagViewItemType ViewType { get { return eBagViewItemType.Item; } }

        public IShowItem item;        

        public ViewItem(IShowItem item)
        {
            this.item = item;
        }
    }

    // 构建用于背包显示的
    public class BagViewData
    {
        List<IBagViewItem> _items = new List<IBagViewItem>();
        public int GetSize()
        {
            return _items.Count;
        }

        public void Add(IBagViewItem item)
        {
            _items.Add(item);            
        }

        public IBagViewItem Get(int index)
        {
            return _items[index];
        }

        public void Sort(Comparison<IBagViewItem> comparison)
        {
            _items.Sort(comparison);
        }
    }
} // namespace Phoenix
