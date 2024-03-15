using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    // 技能
    // 
    public interface ISkill
    {
        // 获取技能所属武器        
        IWeapon GetCurWeapon();
    }

    public interface IWeapon
    {
        int GetMinDmg();
        int GetMaxDmg();
        // 攻速
        float GetSpeed();
        void Reset();
    }

    public interface ICharacter
    {
        // 主武器
        IWeapon GetMainWeapon();
        // 副手武器
        IWeapon GetOffHandWeapon();
    }

    // 技能结算的结果
    public interface IFormulaResult
    {
    }

    // 某个结算公式
    public interface IFormula
    {
        // 执行
        IFormulaResult Apply(ISkill skill, ICharacter src, ICharacter tar);
    }
}// namespace Phoenix
