package monster

import (
	"mmo/modules/csv"
	"mmo/modules/csv/entry"
	"mmo/modules/fight/common"
)

// 怪物属性初始化器

type AttrInitor struct {
	cfgId     string
	md        *entry.Monster
	td        *entry.MonsterTemplate
	character common.ICharacter
}

func NewAttrInitor(cfgId string) *AttrInitor {
	i := &AttrInitor{
		cfgId: cfgId,
	}
	i.md = csv.Monster.GetEntry(i.cfgId)
	if i.md != nil {
		i.td = csv.MonsterTemplate.GetEntry(i.md.Template)
	}
	return i
}

func (i *AttrInitor) OnAdded(character common.ICharacter) {
	i.character = character
	character.GetEvents().SubscribeWithReceiver("onlevelchanged", i, i.onLevelChanged)
}

func (i *AttrInitor) apply(character common.ICharacter, flag int) {
	td := i.td
	if td == nil {
		return
	}

	level := character.GetIntValue(common.Level)

	character.OffsetIntBase(common.HPMax, flag*int(td.HP+td.HPLv*float32(level-1)))
	character.OffsetIntBase(common.PhysicPower, flag*10)
	character.OffsetIntBase(common.PhysicArmor, flag*int(td.Armor+td.ArmorLv*float32(level-1)))
}

func (i *AttrInitor) Equip(character common.ICharacter) {
	i.apply(character, 1)
}

func (i *AttrInitor) Unequip(character common.ICharacter) {
	i.apply(character, -1)
}

func (i *AttrInitor) onLevelChanged(args ...any) {
	td := i.td
	if td == nil {
		return
	}

	oldV := args[0].(int)
	newV := args[1].(int)

	hpChanged := int(float32(newV-oldV) * td.HPLv)
	defendChanged := float32(newV-oldV) * td.ArmorLv

	i.character.OffsetIntBase(common.HPMax, hpChanged)
	i.character.OffsetIntBase(common.PhysicArmor, int(defendChanged))
}
