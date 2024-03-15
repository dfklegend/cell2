package builder

import (
	"mmo/modules/fight/buf"
	charimpls "mmo/modules/fight/character/impls"
	"mmo/modules/fight/common"
	"mmo/modules/fight/equip"
	"mmo/modules/fight/skill"
)

// 统一

type CharBuilder struct {
}

func NewBuilder() *CharBuilder {
	return &CharBuilder{}
}

func (cb *CharBuilder) Build() common.ICharacter {
	b := charimpls.NewBuilder()
	b.WithSkill(skill.NewCtrl())
	b.WithSkillTable(skill.NewSkillTable())
	b.WithBufCtrl(buf.NewCtrl())
	b.WithSlots(equip.NewSlots(10))
	return b.Build()
}
