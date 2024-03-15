namespace Phoenix.Game.FightEmulator.BagSystem
{	
    public static class PlayerBagsBuilder
    {
        public static PlayerBags CreateBags()
        {
            PlayerBags bags = new PlayerBags();
            bags.Init((int)eBagType.MaxBag);
            Bag bag;
            
            bag = CreateBag((int)eBagType.BagTemp, 10);
            bags.SetBag((int)eBagType.BagTemp, bag);            
            bag = CreateBag((int)eBagType.BagItem, (int)eEquipSlot.Max);
            bags.SetBag((int)eBagType.BagItem, bag);
            bag = CreateBag((int)eBagType.BagItem, 100);
            bags.SetBag((int)eBagType.BagItem, bag);
            return bags;
        }

        public static Bag CreateBag(int bagType, int maxSlot)
        {
            Bag bag = new Bag();
            bag.InitBag(bagType, maxSlot);
            return bag;
        }
    }
} // namespace Phoenix
