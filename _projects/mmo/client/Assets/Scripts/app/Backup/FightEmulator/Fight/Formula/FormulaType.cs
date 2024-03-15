using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    public enum eFormulaType
    {
        SimpleTest = 0,         // 测试
        NoDmg,
        PhysicMelee,            // 武器近战
        Heal,                   // 治疗
    }

    public static partial class FormulaDefine
    {
        public static string[] formulaTypes = { "测试",
            "无伤害",
            "物理近战",
            "治疗"
        };
    }
}// namespace Phoenix
