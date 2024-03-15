package charimpls

import (
	"mmo/modules/fight/attr/impls"
	"mmo/modules/fight/common"
)

var (
	minHPMaxClamper = impls.NewMinValueClapmper(1)
)

// Builder 用来build character
type Builder struct {
	attrSize   int
	skill      common.ISkillCtrl
	skillTable common.ISkillTable
	buf        common.IBufCtrl
	slots      common.IEquipSlots
}

func NewBuilder() *Builder {
	return &Builder{}
}

func NewCharacterBuilder() *Builder {
	return NewBuilder()
}

func (b *Builder) WithSkill(ctrl common.ISkillCtrl) {
	b.skill = ctrl
}

func (b *Builder) WithSkillTable(table common.ISkillTable) {
	b.skillTable = table
}

func (b *Builder) WithBufCtrl(ctrl common.IBufCtrl) {
	b.buf = ctrl
}

func (b *Builder) WithSlots(slots common.IEquipSlots) {
	b.slots = slots
}

func (b *Builder) MakeAttrs(c *Character) {
	c.newAttrs(common.MaxAttrNum)

	c.setLevel(1)
	// 设置属性clamper
	c.GetAttr(common.HPMax).SetValueClamper(minHPMaxClamper)
}

func (b *Builder) Build() *Character {
	c := NewCharacter()

	b.MakeAttrs(c)
	c.skill = b.skill
	c.skillTable = b.skillTable
	c.setBufCtrl(b.buf)
	if b.slots != nil {
		c.slots = b.slots
		c.AddEquipGroup(b.slots)
	}
	return c
}
