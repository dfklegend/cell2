package test

import (
	"mmo/modules/fight/buf"
	charimpls "mmo/modules/fight/character/impls"
	"mmo/modules/fight/equip"
	"mmo/modules/fight/skill"
)

func createCharBuilder() *charimpls.Builder {
	b1 := charimpls.NewCharacterBuilder()
	b1.WithSkillTable(skill.NewSkillTable())
	b1.WithSkill(skill.NewCtrl())
	b1.WithBufCtrl(buf.NewCtrl())
	b1.WithSlots(equip.NewSlots(10))
	return b1
}
