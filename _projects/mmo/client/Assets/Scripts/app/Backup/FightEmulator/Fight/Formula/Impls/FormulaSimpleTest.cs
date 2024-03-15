using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    public class FormulaSimpleTest : IFormula
    {
        // 简单物理伤害
        public IFormulaResult Apply(ISkill skill, ICharacter srcChar, ICharacter tarChar)
        {
            if (srcChar == null || tarChar == null)
                return null;
            Character src = srcChar as Character;
            Character tar = tarChar as Character;            

            var result = new FormulaResult();
            int weaponDmg = 10;
            float dmg = (weaponDmg + src.attrs.GetAttr(AttrDefine.MeleePower).intFinal/14 + 1);

            // 护甲免伤
            float armor = tar.attrs.GetAttr(AttrDefine.Armor).final;
            int srcLevel = src.GetLevel();
            float reduce = armor / (armor + 85 * srcLevel + 400);
            dmg = dmg*(1.0f-reduce);

            if (dmg < 0)
                dmg = 0;

            result.data.Dmg = dmg;
            return result;
        }
    }
}// namespace Phoenix
