using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;

namespace Phoenix.Game
{	
    // 界面显示item数据
    // 物品显示的抽象
    public interface IShowItem
    {     
        // 物品id
        string GetItemId();
        int GetItemType();
        // 1: 不能堆叠
        // >1: 可以堆叠
        int GetMaxStack();
        int GetStack();       
        bool CanStack();

        string GetName();
        string GetIcon();
    }

    public interface IShowItemStyle
    {
    }

    public interface IShowItemHandler
    {
        void OnClick(IShowItemStyle style, IShowItem item);
    }
} // namespace Phoenix
