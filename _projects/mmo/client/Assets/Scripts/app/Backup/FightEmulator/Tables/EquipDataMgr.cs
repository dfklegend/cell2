using UnityEngine;
using System.Reflection;
using Phoenix.csv;

namespace Phoenix.Game.FightEmulator
{

    [System.Reflection.ObfuscationAttribute(Exclude = true)]
    public class EquipData
    {
        public string id;

        public int type = 0;                        // 装备类型
        public eEquipBaseType baseType = 0;         // 基础类型，用来确定装备合法性

        public float speed;
        public int minDmg;
        public int maxDmg;
        public int block;
        public string attr0;
        public float value0;
        public string attr1;
        public float value1;
        public string attr2;
        public float value2;
        public string attr3;
        public float value3;

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
