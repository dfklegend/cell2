package attrinitor

import (
	"mmo/modules/fight/attr"
	"mmo/modules/fight/common"
)

// 一个测试的怪物属性初始化器
// hp = 100 + level*10
// attack = 10 + level*1.5
// defend = 1 + level

type TestBaseAttrInitor struct {
	character common.ICharacter
}

func NewTestBaseAttrInitor() *TestBaseAttrInitor {
	return &TestBaseAttrInitor{}
}

func (m *TestBaseAttrInitor) OnAdded(character common.ICharacter) {
	m.character = character
	character.GetEvents().SubscribeWithReceiver("onlevelchanged", m, m.onLevelChanged)
}

func (m *TestBaseAttrInitor) apply(character common.ICharacter, flag int) {
	level := character.GetIntValue(common.Level)

	hp := 100 + level*10
	attack := 10 + float32(level)*1.5
	defend := 1 + level

	character.OffsetBase(common.HPMax, attr.Value(hp))
	character.OffsetPercent(common.HPMax, 0.1)
	character.OffsetIntBase(common.PhysicPower, int(attack)*flag)
	character.OffsetIntBase(common.PhysicArmor, int(defend)*flag)
}

func (m *TestBaseAttrInitor) Equip(character common.ICharacter) {
	m.apply(character, 1)
}

func (m *TestBaseAttrInitor) Unequip(character common.ICharacter) {
	m.apply(character, -1)
}

func (m *TestBaseAttrInitor) onLevelChanged(args ...any) {
	oldV := args[0].(int)
	newV := args[1].(int)

	hpChanged := (newV - oldV) * 10
	attackChanged := float32(newV-oldV) * 1.5
	defendChanged := float32(newV-oldV) * 1

	m.character.OffsetBase(common.HPMax, attr.Value(hpChanged))
	m.character.OffsetIntBase(common.PhysicPower, int(attackChanged))
	m.character.OffsetIntBase(common.PhysicArmor, int(defendChanged))
}
