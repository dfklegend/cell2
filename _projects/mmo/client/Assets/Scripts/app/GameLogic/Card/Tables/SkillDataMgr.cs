using UnityEngine;
using System.Reflection;
using Phoenix.csv;

namespace Phoenix.Game.Card
{
    [System.Reflection.ObfuscationAttribute(Exclude = true)]
    public class SkillData
    {
        public string id;

        public string name;
        public string desc;
        public string icon;
        public Skill.eSkillType type = Skill.eSkillType.Normal;
        public int tarType;
        public int canSelect = 0;
        public int cost = 0;
        public float cd = 0f;

        public string action;
        public int normalAttack;
        public float totalTime = 0;
        public float hitTime = 0;
        public int formulaType;
        public float baseDmg = 0;
        public float baseDmgLv = 0;
        public float powerMultiplier = 0;
        public float powerMultiplierLv = 0;

        // 
        public string[] subSkills;

        public string Key() { return id; }
    }


    public class SkillDataMgr : TableManager<string, SkillData>
    {
        static SkillDataMgr sInstance = null;
        public static SkillDataMgr It
        {
            get
            {
                if (sInstance == null)
                    sInstance = new SkillDataMgr();
                return sInstance;
            }
        }

        public override string MakeKey(SkillData obj) { return obj.Key(); }
        public override int indexRead() { return 2; }
        public override string TablePath() { return "csv/skill"; }
        
    }
}
