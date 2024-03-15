using System;
using System.Collections.Generic;
using UnityEngine;
using Phoenix.Core;
using Phoenix.Entity;


namespace Phoenix.Game.FightEmulator
{
    // 主手武器，近战
    
    public class WeaponMelee : IFormula
    {
        // 简单物理伤害
        public IFormulaResult Apply(ISkill iskill, ICharacter srcChar, ICharacter tarChar)
        {
            if (srcChar == null || tarChar == null)
                return null;
            Character src = srcChar as Character;
            Character tar = tarChar as Character;
            Skill.Skill skill = iskill as Skill.Skill;

            var skillCfg = skill.cfgData;
            var skillLv = 1;

            var result = new FormulaResult();


            result.data.hand = skill.GetHandType();

            // 计算命中            
            float hitChance = 1f;
            
            float weaponBaseMiss = 0.05f;
            if (skill.GetHandType() == Skill.eHandType.OffHand)
                weaponBaseMiss = 0.24f;
            
            hitChance -= weaponBaseMiss;

            float hitAttrFactor = FormulaUtil.CalcHitChangeFromPoint(
                src.attrs.GetAttr(AttrDefine.PhysicHit).intFinal, tar.level);
            float levelFactor = 0.03f * (src.level - tar.level);
            hitChance += hitAttrFactor + levelFactor;
            if(MathUtil.RandomF(0,1) > hitChance)
            {
                result.data.hit = false;
                return result;
            }

            // TODO:招架

            result.data.hit = true;            
            var weapon = skill.GetCurWeapon();
            int weaponDmg = CalcWeaponDmg(weapon);
            // apply武器伤害增强
            weaponDmg = (int)((float)weaponDmg * (1f + src.attrs.GetAttr(AttrDefine.WeaponEnhance).final));

            // 技能基础伤害 + 武器伤害*武器伤害比 + 攻强加成
            float weaponFactor = skillCfg.powerMultiplier + skillCfg.powerMultiplierLv * skillLv;
            
            // 计算攻强加成
            float baseDmg = skillCfg.baseDmg + skillCfg.baseDmgLv*skillLv +
                weaponDmg* weaponFactor + src.attrs.GetAttr(AttrDefine.MeleePower).intFinal/14 + 1;                


            // ----- 伤害减免
            // 护甲免伤
            float armor = tar.attrs.GetAttr(AttrDefine.Armor).final;
            int srcLevel = src.GetLevel();
            float reduce = armor / (armor + 85 * srcLevel + 400);

            Log.LogCenter.Default.Debug($"护甲: {armor} 免伤: {reduce}");
            var finalDmg = baseDmg*(1.0f-reduce);

            // 盾牌格挡
            var shieldBlock = tar.GetShieldBlock();
            var blockChance = 0.05f +
                FormulaUtil.StandardPointToChance(
                tar.attrs.GetAttr(AttrDefine.BlockChance).intFinal, src.level);
            if (shieldBlock > 0 && MathUtil.HitChance(blockChance) )
            {
                result.data.block = true;
                result.data.dmgBlocked = shieldBlock;
                finalDmg -= result.data.dmgBlocked;
                if (finalDmg < 1)
                    finalDmg = 1;
            }

                // 计算暴击
            float criticalChance = 0.05f;
            // 属性
            criticalChance += FormulaUtil.StandardPointToChance(
                src.attrs.GetAttr(AttrDefine.PhysicCritical).intFinal, tar.level);
            // TODO: 韧性
            if (MathUtil.HitChance(criticalChance))
            {
                finalDmg *= 2f;
                result.data.critical = true;
            }

            result.data.Dmg = finalDmg;
            return result;
        }

        private int CalcWeaponDmg(IWeapon weapon)
        {
            if (weapon == null)
                return 0;
            var dmg = MathUtil.RandomI(weapon.GetMinDmg(), weapon.GetMaxDmg() + 1);
            Log.LogCenter.Default.Debug($"武器伤害: {dmg} {weapon.GetMinDmg()}-{weapon.GetMaxDmg()}");
            return dmg;
        }
    }
}// namespace Phoenix
