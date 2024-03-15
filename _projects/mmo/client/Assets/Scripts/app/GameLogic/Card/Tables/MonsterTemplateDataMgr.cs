using UnityEngine;
using System.Reflection;
using Phoenix.csv;

namespace Phoenix.Game.Card
{

    [ObfuscationAttribute(Exclude = true)]
    public class MonsterTemplateData
    {
        public string id;

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


    public class MonsterTemplateDataMgr : TableManager<string, MonsterTemplateData>
    {
        static MonsterTemplateDataMgr sInstance = null;
        public static MonsterTemplateDataMgr It
        {
            get
            {
                if (sInstance == null)
                    sInstance = new MonsterTemplateDataMgr();
                return sInstance;
            }
        }

        public override string MakeKey(MonsterTemplateData obj) { return obj.Key(); }
        public override int indexRead() { return 2; }
        public override string TablePath() { return "csv/monster_template"; }
        
    }
}
