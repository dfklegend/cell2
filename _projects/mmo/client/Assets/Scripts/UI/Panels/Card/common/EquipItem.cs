using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Log;

namespace Phoenix.Game.Card
{
    public class EquipItem : IShowItem
    {
        private string _id;
        private ItemData _item;
        private EquipData _equip;
        private string _iconPath = "";

        public EquipItem(string equipId)
        {
            _id = equipId;
            _item = ItemDataMgr.It.GetItem(_id);
            _equip = EquipDataMgr.It.GetItem(_id);
            if(_item != null )
            {
                _iconPath = $"items/{_item.icon}";
            }            
        }

        public bool CanStack()
        {
            return false;
        }

        public string GetIcon()
        {           
            return _iconPath;
        }

        public string GetItemId()
        {
            return _id;
        }

        public int GetItemType()
        {
            return 0;
        }

        public int GetMaxStack()
        {
            return 1;
        }

        public string GetName()
        {
            return _id;
        }

        public int GetStack()
        {
            return 1;
        }       
    }

} // namespace Phoenix
