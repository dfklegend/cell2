
using System.Collections.Generic;
using Phoenix.Core;
using Phoenix.Entity;

namespace Phoenix.Game.FightEmulator.BagSystem
{	    
    // 移除了之后,背包中间会留空，保证index不变
    public class Bag : IBag
    {
        private int _bagType = -1;
        private int _maxSlotNum = 0;
        private List<IBagItem> _items =
            new List<IBagItem>();
        private List<IBagFeature> _features = 
            new List<IBagFeature>();

        public void InitBag(int bagType, int maxSlot)
        {
            _bagType = bagType;
            _maxSlotNum = maxSlot;
            extendSlots(maxSlot);
        }

        private void extendSlots(int size)
        {
            if (_items.Count >= size)
                return;
            int begin = _items.Count;
            for(var i = begin; i < size; i ++ )
            {
                _items.Add(null);
            }
        }

        public int GetRealItemNum()
        {
            int num = 0;
            for (var i = 0; i < _items.Count; i++)
            {
                var item = _items[i];
                if (item == null)
                    continue;
                num++;
            }
            return num;
        }

        public int GetActiveNum()
        {
            return _items.Count;
        }

        public int GetMaxSlotNum()
        {
            return _items.Count;
        }

        public int GetFreeNum()
        {
            return GetMaxSlotNum() - GetRealItemNum();
        }

        public void AddFeature(IBagFeature feature)
        {
            _features.Add(feature);
        }

        public bool IsValidSlot(int index)
        {
            if (index < 0 || index >= _items.Count)
                return false;
            return true;
        }

        public IBagItem GetItem(int index)
        {
            if (index < 0 || index >= _items.Count)
                return null;
            return _items[index];
        }

        public T GetItem<T>(int index)
            where T: class, IBagItem
        {
            var v = GetItem(index);
            return v as T;
        }

        public bool SetItem(int index, IBagItem item)
        {
            if (!IsValidSlot(index))
                return false;
            if (!checkCanSet(index, item))
                return false;
            // 背包上有物品
            if (GetItem(index) != null)
                return false;

            Log.LogCenter.Default.Debug("add item {0} in {1}", item.GetItemId(), index);

            item.SetIndex(_bagType, index);
            _items[index] = item;
            return true;
        }

        public bool AddItem(IBagItem item)
        {
            // 物品加入背包
            if (item.CanStack())
                return addStackableItem(item);
            return addNonStackableItem(item);
        }

        public bool SetOrAdd(int index, IBagItem item)
        {
            if(index != -1)
            {
                if (!IsValidSlot(index))
                    return false;
                if (SetItem(index, item))
                    return true;
            }            
            return AddItem(item);
        }

        private bool addStackableItem(IBagItem item)
        {
            return false;
        }       
        
        private bool addNonStackableItem(IBagItem item)
        {
            var index = findEmptySlot();
            if (index == -1)
                return false;
            return SetItem(index, item);
        }

        private int findEmptySlot()
        {
            for(var i = 0; i < _items.Count; i ++)
            {
                if (_items[i] == null)
                    return i;
            }
            return -1;
        }

        public bool RemoveItem(int index)
        {
            if (!IsValidSlot(index))
                return false;

            var item = _items[index];
            if (item == null)
                return false;

            Log.LogCenter.Default.Debug("remove item {0} from {1}", item.GetItemId(), index);

            item.SetIndex(-1, -1);
            _items[index] = null;
            
            return true;
        }       

        private bool checkCanSet(int index, IBagItem item)
        {
            if (_features.Count == 0)
                return true;
            for (var i = 0; i < _features.Count; i++)
                if (!_features[i].CanSet(this, index, item))
                    return false;
            return true;
        }

        public void DumpInfo()
        {
            Log.LogCenter.Default.Debug("bag:");
            for (var i = 0; i < _items.Count; i++)
            {
                var item = _items[i];
                if (item == null)
                {
                    Log.LogCenter.Default.Debug("  {0} null", i);
                    continue;
                }
                    
                Log.LogCenter.Default.Debug("  {0} {1} {2}",
                    i, item.GetItemId(), item.GetIndex());
            }            
        }
    }
} // namespace Phoenix
