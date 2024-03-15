package test

import (
	"github.com/dfklegend/cell2/utils/event/light"

	"mmo/modules/fight/common"
	"mmo/modules/fight/skill/define"
)

//	额外增加10点物强
type BonusPhysicPower struct {
}

func (b *BonusPhysicPower) Start(events *light.EventCenter) {
	b.bindEvents(true, events)
}

func (b *BonusPhysicPower) bindEvents(bind bool, events *light.EventCenter) {
	light.BindEventWithReceiver(bind, events, define.FormulaEventSrcBeforeCalcDmg1, b, b.onFormula)
}

func (b *BonusPhysicPower) onFormula(args ...any) {
	dmg := args[0].(*common.DmgInstance)
	dmg.BonusPhysicPower += 10
}
