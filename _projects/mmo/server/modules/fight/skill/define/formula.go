package define

import (
	"mmo/modules/fight/common"
)

//	无状态

type IFormula interface {
	Apply(skill ISkill, src common.ICharacter, tar common.ICharacter, data *common.DmgInstance)
}

const (
	// 	FormulaEventSrcBeforeCalcDmg1 计算dmg前触发，可以提供一些额外物理法术加成
	FormulaEventSrcBeforeCalcDmg1 = "发起结算开始1"
	// 	FormulaEventTarBeforeCalcDmg1 计算dmg前触发，可以提供一些额外物理法术护甲
	FormulaEventTarBeforeCalcDmg1 = "受到结算开始1"

	FormulaEventSrcAfterCalcDmg2 = "发起结算开始2"
	FormulaEventTarAfterCalcDmg2 = "受到结算开始2"

	// 	FormulaEventProcessAbsorb 处理伤害吸收，一般由buf处理此事件
	FormulaEventProcessAbsorb = "处理吸收"

	FormulaEventSrcAfterApplyDmg3 = "伤害命中"
	FormulaEventTarAfterApplyDmg3 = "受到伤害"
)
