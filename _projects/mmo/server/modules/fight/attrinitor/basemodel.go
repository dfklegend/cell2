package attrinitor

import (
	"mmo/modules/fight/attr"
	"mmo/modules/fight/common"
)

// 一个测试的

type BaseModelInitor struct {
	character common.ICharacter
}

func NewBaseModelInitor() *BaseModelInitor {
	return &BaseModelInitor{}
}

func (m *BaseModelInitor) OnAdded(character common.ICharacter) {
	m.character = character
}

func (m *BaseModelInitor) apply(c common.ICharacter, flag int) {

	hp := 10000
	attack := 100
	defend := 50
	//
	attackSpeed := 0.5

	c.OffsetBase(common.HPMax, attr.Value(hp))
	c.SetBaseValue(common.AttackSpeed, attackSpeed)
	c.OffsetIntBase(common.PhysicPower, int(attack)*flag)
	c.OffsetIntBase(common.PhysicArmor, int(defend)*flag)
}

func (m *BaseModelInitor) Equip(character common.ICharacter) {
	m.apply(character, 1)
}

func (m *BaseModelInitor) Unequip(character common.ICharacter) {
	m.apply(character, -1)
}
