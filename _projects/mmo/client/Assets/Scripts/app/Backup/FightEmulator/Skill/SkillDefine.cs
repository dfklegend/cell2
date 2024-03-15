using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;

using Phoenix.Game.FightEmulator;


namespace Phoenix.Game.Skill
{	    
    /*
     * [channel]可选
     * prehit   play action
     * (hitpoint)
     * apply cost
     * posthit
     */
    public enum eSkillPhase
    {
        Init = 0,
        Channel,        // 吟唱
        Prehit,         // 前摇
        Posthit,        // 后摇
        Over,           // 结束
        SubskillRunning // 子技能执行中
    }
    /*
     * 如果有子技能，主技能命中点后，就开始启动子技能
     */



    public enum eSkillType
    {
        Normal = 0,
        TargetBullet,    // 对目标发射子弹
    }

    public static partial class SkillDefine
    {
        public const int MaxSkillSlot = 5;
        public const int IndexSuperSkill = MaxSkillSlot - 1;

        public const string MeleeSkillId = "普攻";
        public const string RangeSkillId = "远程普攻";
    }

} // namespace Phoenix
