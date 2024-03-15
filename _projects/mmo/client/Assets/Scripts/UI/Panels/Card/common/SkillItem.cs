using Phoenix.Utils;
using UnityEngine;
using UnityEngine.UI;
using Phoenix.Log;

namespace Phoenix.Game.Card
{
    // 用于技能类显示
    public class SkillItem : IShowItem
    {
        private string _id;
        private SkillData _item;        
        private string _iconPath;

        public SkillItem(string id)
        {
            _id = id;
            _item = SkillDataMgr.It.GetItem(_id);            
            if(_item != null )
            {
                _iconPath = $"skills/{_item.icon}";
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
