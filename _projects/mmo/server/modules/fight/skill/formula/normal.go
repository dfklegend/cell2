package formula

import (
	"math/rand"

	l "github.com/dfklegend/cell2/utils/logger"

	"mmo/modules/csv/entry"
	"mmo/modules/fight/attr"
	"mmo/modules/fight/common"
	"mmo/modules/fight/common/skilldef"
	define2 "mmo/modules/fight/skill/define"
)

/*
事件列表

	所有者: 攻击方 	事件名: 发起结算开始1
		可以增加本次额外攻击，额外穿透
	所有者: 目标 		事件名: 受到结算开始1
		可以增加本次额外防御参数

	...
	计算攻防
	得到伤害值
	...


	所有者: 攻击方 	事件名: 发起结算开始2
		可以计算额外伤害，额外增伤
	所有者: 目标 		事件名: 受到结算开始2
		可以计算额外减伤

	所有者: 目标 		事件名: 处理吸收
		处理盾吸收伤害

	...
	applyDmg
	...

	所有者: 攻击方 	事件名: 伤害命中
		可以做一些命中事件处理
	所有者: 目标 		事件名: 受到伤害



*/

type FuncCalcArmor func(src common.ICharacter, tar common.ICharacter, data *common.DmgInstance) attr.Value

func init() {
	Register(skilldef.FromulaNormal, normal)
}

var normal = &Normal{}

type Normal struct {
}

func (n *Normal) Apply(skill define2.ISkill, src common.ICharacter, tar common.ICharacter, data *common.DmgInstance) {
	cfg := skill.GetCfg()

	data.Src = src
	data.Tar = tar
	data.SrcProxy = src.GetProxy()
	data.TarProxy = tar.GetProxy()
	data.DmgType = cfg.DmgType.Type

	srcEvents := src.GetEvents()
	tarEvents := tar.GetEvents()

	srcEvents.Publish(define2.FormulaEventSrcBeforeCalcDmg1, data)
	tarEvents.Publish(define2.FormulaEventTarBeforeCalcDmg1, data)

	// TODO: 处理命中

	// 处理基础伤害
	var dmg float32
	switch data.DmgType {
	case skilldef.DmgTypePhysic:
		dmg = n.calcPhysicDmg(skill, src, tar, data)
	case skilldef.DmgTypeMagic:
		dmg = n.calcMagicDmg(skill, src, tar, data)
	case skilldef.DmgTypeReal:
		dmg = n.calcRealDmg(skill, src, tar, data)
	}
	if dmg < 0 {
		dmg = 0
	}

	//// 概率暴击
	//critical := false
	//criticalRate := float32(0.3)
	//if rand.Float32() < criticalRate {
	//	critical = true
	//	dmg *= 2
	//}

	critical, finalDmg := calcCritical(cfg, src, tar, data, dmg)
	data.Critical = critical
	data.Dmg = int(finalDmg)

	srcEvents.Publish(define2.FormulaEventSrcAfterCalcDmg2, data)
	tarEvents.Publish(define2.FormulaEventTarAfterCalcDmg2, data)

	if data.TarBonusDmgReduceP > 0 {
		l.L.Infof("%v", data.TarBonusDmgReduceP)
	}

	// 伤害增伤，减伤处理
	dmgAfterModify := calcDmgModify(float64(finalDmg), data.BonusDmgEnhanceP, data.TarBonusDmgReduceP,
		data.BonusDmgEnhance, data.TarBonusDmgReduce)
	data.Dmg = int(dmgAfterModify)

	// 可以获取技能单独增伤，技能种类增伤
	// 获取目标对某个，某类技能的减伤

	// 盾吸收伤害
	tarEvents.Publish(define2.FormulaEventProcessAbsorb, data)

	// 实际产生伤害
	dmgApplied := tar.ApplyDmg(data.Dmg)
	if dmgApplied == 0 {
		data.NoDmg = true
	}

	if dmgApplied > 0 && tar.IsDead() {
		data.TarBeKilled = true
	}

	srcEvents.Publish(define2.FormulaEventSrcAfterApplyDmg3, data)
	tarEvents.Publish(define2.FormulaEventTarAfterApplyDmg3, data)
}

func (n *Normal) calcPhysicDmg(skill define2.ISkill, src common.ICharacter, tar common.ICharacter, data *common.DmgInstance) float32 {
	return n.calcDmg(skill, src, tar, data, calcPhysicArmor)
}

func (n *Normal) calcMagicDmg(skill define2.ISkill, src common.ICharacter, tar common.ICharacter, data *common.DmgInstance) float32 {
	return n.calcDmg(skill, src, tar, data, calcMagicArmor)
}

func (n *Normal) calcRealDmg(skill define2.ISkill, src common.ICharacter, tar common.ICharacter, data *common.DmgInstance) float32 {
	return n.calcDmg(skill, src, tar, data, calcRealArmor)
}

func (n *Normal) calcDmg(skill define2.ISkill, src common.ICharacter, tar common.ICharacter, data *common.DmgInstance, fCalcArmor FuncCalcArmor) float32 {
	cfg := skill.GetCfg()
	level := float32(1)
	baseDmg := float64(cfg.BaseDmg + cfg.BaseDmgLv*1)

	adFactor := float64(cfg.AD + cfg.ADLv*level)
	physic := src.GetValue(common.PhysicPower)
	apFactor := float64(cfg.AP + cfg.APLv*level)
	magic := src.GetValue(common.MagicPower)

	// 基础伤害加上 物理加成 + 法伤加成
	damage := baseDmg + adFactor*(physic+data.BonusPhysicPower) + apFactor*(magic+data.BonusMagicPower)

	armor := fCalcArmor(src, tar, data)
	// 伤害公式
	// 减法
	damage -= armor
	return float32(damage)
}

func calcPhysicArmor(src common.ICharacter, tar common.ICharacter, data *common.DmgInstance) attr.Value {
	// 根据伤害类型，来选择护甲
	// 计算目标护甲 (带入穿透)
	armor := tar.GetValue(common.PhysicArmor) + data.TarBonusArmor

	// 百分比护甲穿透
	armor = armor * (1 - data.BonusPenetrateP)

	// 数值穿透
	armor -= data.BonusPenetrate
	if armor < 0 {
		armor = 0
	}
	return armor
}

func calcMagicArmor(src common.ICharacter, tar common.ICharacter, data *common.DmgInstance) attr.Value {
	// 根据伤害类型，来选择护甲
	// 计算目标护甲 (带入穿透)
	armor := tar.GetValue(common.MagicArmor) + data.TarBonusMagicArmor

	// 百分比护甲穿透
	armor = armor * (1 - data.BonusMagicPenetrateP)

	// 数值穿透
	armor -= data.BonusMagicPenetrate
	if armor < 0 {
		armor = 0
	}
	return armor
}

func calcRealArmor(src common.ICharacter, tar common.ICharacter, data *common.DmgInstance) attr.Value {
	return 0
}

// calcCritical 只有普攻才会暴击
func calcCritical(cfg *entry.Skill, src common.ICharacter, tar common.ICharacter, data *common.DmgInstance, dmg float32) (critical bool, finalDmg float32) {
	if !cfg.NormalAttack {
		return false, dmg
	}

	// 初始暴击
	initCriticalRate := 0.05
	criticalRate := src.GetValue(common.CriticalRate) + initCriticalRate
	criticalRate += data.BonusCriticalRate
	// 暴击抵抗
	criticalRate -= data.BonusCriticalDmgResist

	critical = rand.Float64() < criticalRate
	critialFactor := 2.0
	critialFactor += data.BonusCriticalDmgRate
	finalDmg = dmg
	if critical {
		finalDmg *= float32(critialFactor)
	}
	return critical, finalDmg
}

func calcDmgModify(dmg float64, dmgEnhanceP, dmgReduceP, dmgEnhance, dmgReduce float64) float64 {
	var factor float64 = 1
	if dmgEnhanceP != 0 || dmgReduceP != 0 {
		enhancefactor := 1 + dmgEnhanceP
		reduceFactor := 1 - dmgReduceP
		// 最小减伤到10%
		if reduceFactor < 0 {
			reduceFactor = 0.1
		}
		factor = enhancefactor * reduceFactor
	}

	curDmg := dmg * factor
	curDmg = curDmg + dmgEnhance - dmgReduce
	if curDmg < 0 {
		curDmg = 0
	}
	return curDmg * factor
}
