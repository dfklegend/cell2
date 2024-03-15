using UnityEngine;
using System.Reflection;
using Phoenix.csv;

namespace Phoenix.Game.Card
{

    [System.Reflection.ObfuscationAttribute(Exclude = true)]
    public class TestData
    {
        public int id;
        public string name;         // 名字        

        public int Key() { return id; }
    }


    public class TestDataMgr : TableManager<int, TestData>
    {
        static TestDataMgr sInstance = null;
        public static TestDataMgr It
        {
            get
            {
                if (sInstance == null)
                    sInstance = new TestDataMgr();
                return sInstance;
            }
        }

        public override int MakeKey(TestData obj) { return obj.Key(); }
        public override int indexRead() { return 2; }
        public override string TablePath() { return "csv/test"; }
        
    }
}
