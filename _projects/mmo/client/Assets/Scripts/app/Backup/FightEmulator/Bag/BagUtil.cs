using System.Collections.Generic;

namespace Phoenix.Game.FightEmulator.BagSystem
{	
    public static class BagUtil
    {
        public static bool RemoveItem(PlayerBags bags, int bagType, int index)
        {
            bags.GetBag(bagType).RemoveItem(index);
            return true;
        }

        public static bool RemoveItem(PlayerBags bags, IBagItem item)
        {
            bags.GetBag(item.GetBagType()).RemoveItem(item.GetIndex());
            return true;
        }
    }
} // namespace Phoenix
