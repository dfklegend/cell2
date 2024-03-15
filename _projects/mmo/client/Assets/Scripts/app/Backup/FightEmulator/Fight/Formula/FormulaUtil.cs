using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    public static class FormulaUtil
    {
        public static IFormulaResult SimpleFormula(ISkill skill, Character src, Character tar)
        {
            return FormulaFactory.It.GetFormula((int)eFormulaType.SimpleTest).Apply(skill, src, tar);
        }

        public static IFormulaResult WeaponFormula(ISkill skill, Character src, Character tar)
        {
            return FormulaFactory.It.GetFormula((int)eFormulaType.PhysicMelee).Apply(skill, src, tar);
        }

        public static IFormulaResult Formula(ISkill skill, Character src, Character tar)
        {
            Skill.Skill obj = skill as Skill.Skill;
            IFormula formula = FormulaFactory.It.GetFormula((int)obj.cfgData.formulaType);
            return formula.Apply(skill, src, tar);            
        }

        public static float StandardPointToChance(int point, int tarLevel)
        {
            // 简单公式
            // 按等级1%命中率需求命中点数
            // 10: 3.5 
            // 30: 6.5 
            // 50: 9.5
            // 60: 11
            // 70: 12.5
            float levelReq = 2f + 0.15f * tarLevel;
            return 0.01f * point / levelReq;
        }

        public static float CalcHitChangeFromPoint(int point, int tarLevel)
        {
            return StandardPointToChance(point, tarLevel);
        }

        public static float CalcCriticalChangeFromPoint(int point, int tarLevel)
        {
            return StandardPointToChance(point, tarLevel);
        }
    }
}// namespace Phoenix
