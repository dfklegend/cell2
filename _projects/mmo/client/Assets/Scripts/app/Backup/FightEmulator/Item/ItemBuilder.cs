using System.Collections.Generic;
using Phoenix.Core;
using Phoenix.Entity;

using Phoenix.Game.FightEmulator;

namespace Phoenix.Game.FightEmulator.ItemSystem
{	
    public static class ItemBuilder 
    {
        public static Item CreateItem(string cfgId)
        {
            var cfg = ItemDataMgr.It.GetItem(cfgId);
            if (cfg == null)
                return null;
            Item item = new Item();
            item.Init(cfgId);
            return item;
        }

        public static EquipItem CreateEquip(string cfgId)
        {
            var cfg = ItemDataMgr.It.GetItem(cfgId);
            if (cfg == null)
                return null;
            EquipItem item = new EquipItem();
            item.Init(cfgId);
            return item;
        }
    }    
} // namespace Phoenix
