using UnityEngine;
using System.Reflection;
using Phoenix.csv;

namespace Phoenix.Game.Card
{

    [System.Reflection.ObfuscationAttribute(Exclude = true)]
    public class EquipData
    {
        public string id;

        public int type = 0;                        // 装备类型
        public int baseType = 0;         // 基础类型，用来确定装备合法性

        public float speed;
        public int minDmg;
        public int maxDmg;
        public int block;
        public AttrType attr0;
        public float v0;
        public AttrType attr1;
        public float v1;
        public AttrType attr2;
        public float v2;
        public AttrType attr3;
        public float v3;

        public string Key() { return id; }
    }


    public class EquipDataMgr : TableManager<string, EquipData>
    {
        static EquipDataMgr sInstance = null;
        public static EquipDataMgr It
        {
            get
            {
                if (sInstance == null)
                    sInstance = new EquipDataMgr();
                return sInstance;
            }
        }

        public override string MakeKey(EquipData obj) { return obj.Key(); }
        public override int indexRead() { return 2; }
        public override string TablePath() { return "csv/equip"; }
        
    }
}
