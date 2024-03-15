using System.Collections.Generic;
using Phoenix.Core;
using Phoenix.Entity;

using Phoenix.Game.FightEmulator;

namespace Phoenix.Game.FightEmulator.ItemSystem
{	
    public class BaseItem : BagSystem.IBagItem
    {
        protected ItemData _cfg;
        public ItemData cfg { get { return _cfg; } }
        protected string _cfgId;

        protected int _stack = 0;
        protected int _bagType = -1;
        protected int _index = -1;

        public virtual void Init(string cfgId)
        {
            _cfg = ItemDataMgr.It.GetItem(cfgId);
            _cfgId = cfgId;
        }

        public void SetIndex(int bagType, int n)
        {
            _bagType = bagType;
            _index = n;
        }

        public int GetIndex()
        {
            return _index;
        }

        public int GetBagType()
        {
            return _bagType;
        }

        public string Name { get { return _cfg.name; } }
        public string Icon { get { return _cfg.icon; } }

        public string GetItemId()
        {
            return _cfgId;
        }

        public int GetItemType()
        {
            return 0;
        }

        public int GetMaxStack()
        {
            return 1;
        }

        public int GetStack()
        {
            return 1;
        }

        public void SetStack(int stack)
        {
        }

        public bool CanStack()
        {
            return GetMaxStack() > 1;
        }
    }

    public class Item: BaseItem
    {
    }    

    // EquipItem

    public class ItemBox : BaseItem
    {
    }
} // namespace Phoenix
