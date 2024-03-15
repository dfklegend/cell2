using UnityEngine;
using System.Reflection;
using Phoenix.csv;

namespace Phoenix.Game.FightEmulator
{

    [System.Reflection.ObfuscationAttribute(Exclude = true)]
    public class BulletData
    {
        public int id;
        public string name;         // 名字        

        public int Key() { return id; }
    }


    public class BulletDataMgr : TableManager<int, BulletData>
    {
        static BulletDataMgr sInstance = null;
        public static BulletDataMgr It
        {
            get
            {
                if (sInstance == null)
                    sInstance = new BulletDataMgr();
                return sInstance;
            }
        }

        public override int MakeKey(BulletData obj) { return obj.Key(); }
        public override int indexRead() { return 2; }
        public override string TablePath() { return "csv/bullet"; }
        
    }
}
