package player

import (
	"mmo/modules/fight/common"
)

const (
	HP            = 1000
	HPLv          = 100
	PhysicPower   = 20
	PhysicPowerLv = 10
	PhysicArmor   = 10
	PhysicArmorLv = 2
)

// 玩家属性初始化器

type AttrInitor struct {
	character common.ICharacter
}

func NewAttrInitor() *AttrInitor {
	return &AttrInitor{}
}

func (i *AttrInitor) OnAdded(character common.ICharacter) {
	i.character = character
	character.GetEvents().SubscribeWithReceiver("onlevelchanged", i, i.onLevelChanged)
}

func (i *AttrInitor) apply(character common.ICharacter, flag int) {
	level := character.GetIntValue(common.Level)

	character.OffsetIntBase(common.HPMax, flag*int(HP+HPLv*float32(level-1)))
	character.OffsetIntBase(common.PhysicPower, flag*(int(PhysicPower+PhysicPowerLv*float32(level-1))))
	character.OffsetIntBase(common.PhysicArmor, flag*int(PhysicArmor+PhysicArmorLv*float32(level-1)))
}

func (i *AttrInitor) Equip(character common.ICharacter) {
	i.apply(character, 1)
}

func (i *AttrInitor) Unequip(character common.ICharacter) {
	i.apply(character, -1)
}

func (i *AttrInitor) onLevelChanged(args ...any) {
	oldV := args[0].(int)
	newV := args[1].(int)

	hpChanged := int(float32(newV-oldV) * HPLv)
	defendChanged := float32(newV-oldV) * PhysicArmorLv
	powerChanged := float32(newV-oldV) * PhysicPowerLv

	i.character.OffsetIntBase(common.HPMax, hpChanged)
	i.character.OffsetIntBase(common.PhysicArmor, int(defendChanged))
	i.character.OffsetIntBase(common.PhysicPower, int(powerChanged))
}
