using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;

namespace Phoenix.Game.FightEmulator.BagSystem
{	
    public interface IBagItem
    {
        int GetIndex();
        void SetIndex(int bagType, int index);
        int GetBagType();

        // 物品id
        string GetItemId();
        int GetItemType();
        // 1: 不能堆叠
        // >1: 可以堆叠
        int GetMaxStack();
        int GetStack();
        void SetStack(int stack);
        bool CanStack();
    }

    public interface IBag
    {
    }

    // 抽象一些背包特性
    // 比如 某些背包只能放某些物品
    public interface IBagFeature
    {
        bool CanSet(IBag bag, int index, IBagItem item);
    }
} // namespace Phoenix
