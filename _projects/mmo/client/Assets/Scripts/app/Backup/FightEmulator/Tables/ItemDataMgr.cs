using UnityEngine;
using System.Reflection;
using Phoenix.csv;

namespace Phoenix.Game.FightEmulator
{

    [System.Reflection.ObfuscationAttribute(Exclude = true)]
    public class ItemData
    {
        public string id;
        public string name;
        public string icon;
        public int type;                // 物品类型
        public int maxStack = 1;

        public string Key() { return id; }
    }


    public class ItemDataMgr : TableManager<string, ItemData>
    {
        static ItemDataMgr sInstance = null;
        public static ItemDataMgr It
        {
            get
            {
                if (sInstance == null)
                    sInstance = new ItemDataMgr();
                return sInstance;
            }
        }

        public override string MakeKey(ItemData obj) { return obj.Key(); }
        public override int indexRead() { return 2; }
        public override string TablePath() { return "csv/item"; }
        
    }
}
