using UnityEngine;
using System.Reflection;
using Phoenix.csv;

namespace Phoenix.Game.Card
{

    [ObfuscationAttribute(Exclude = true)]
    public class MonsterData
    {
        public string id;
        public string name;
        public string template;
        public int level;
        public string icon;

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
