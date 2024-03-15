using UnityEngine;
using System.Reflection;
using Phoenix.csv;

namespace Phoenix.Game.FightEmulator
{

    [System.Reflection.ObfuscationAttribute(Exclude = true)]
    public class MonsterData
    {
        public string id;
        public string name;
        public int baseLevel;
        public int level;
        public float weaponSpeed;
        public float guardRange = 6.0f;
        public float attackRange = 2.0f;
        public float minDmg;
        public float minDmgLv;
        public float maxDmg;
        public float maxDmgLv;
        public float armor;
        public float armorLv;
        public float hp;
        public float hpLv;

        public string Key() { return id; }
    }


    public class MonsterDataMgr : TableManager<string, MonsterData>
    {
        static MonsterDataMgr sInstance = null;
        public static MonsterDataMgr It
        {
            get
            {
                if (sInstance == null)
                    sInstance = new MonsterDataMgr();
                return sInstance;
            }
        }

        public override string MakeKey(MonsterData obj) { return obj.Key(); }
        public override int indexRead() { return 2; }
        public override string TablePath() { return "csv/monster"; }
        
    }
}
