package charimpls

import (
	"log"
	"testing"

	"github.com/dfklegend/cell2/utils/event/light"
	"github.com/stretchr/testify/assert"

	"mmo/modules/fight/attr"
	"mmo/modules/fight/common"
)

// HPMax = 100 + level*10
type MonsterTestAttr struct {
	character common.ICharacter
}

func (m *MonsterTestAttr) OnAdded(character common.ICharacter) {
	m.character = character
	character.GetEvents().SubscribeWithReceiver("onlevelchanged", m, m.onLevelChanged)
}

func (m *MonsterTestAttr) Equip(character common.ICharacter) {
	level := character.GetIntValue(common.Level)

	hp := 100 + level*10

	character.OffsetBase(common.HPMax, attr.Value(hp))
	character.OffsetPercent(common.HPMax, 0.1)
}

func (m *MonsterTestAttr) Unequip(character common.ICharacter) {
	level := character.GetIntValue(common.Level)

	hp := 100 + level*10
	character.OffsetBase(common.HPMax, attr.Value(-hp))
	character.OffsetPercent(common.HPMax, -0.1)
}

func (m *MonsterTestAttr) onLevelChanged(args ...any) {
	oldV := args[0].(int)
	newV := args[1].(int)

	hpChanged := (newV - oldV) * 10
	m.character.OffsetBase(common.HPMax, attr.Value(hpChanged))
}

// 验证属性构建器
func TestBase(t *testing.T) {
	b := NewBuilder()
	c := b.Build()
	c.events = light.NewEventCenter()

	m := &MonsterTestAttr{}

	c.ChangeLevel(1)

	c.AddEquipGroup(m)
	log.Printf("%v\n", c.findEquipGroupIndex(m))

	log.Printf("%v\n", c.GetAttr(common.HPMax).GetIntValue())
	assert.Equal(t, 121, c.GetAttr(common.HPMax).GetIntValue())

	c.ChangeLevel(10)
	log.Printf("%v\n", c.GetAttr(common.HPMax).GetIntValue())
	assert.Equal(t, 220, c.GetAttr(common.HPMax).GetIntValue())

	m.Unequip(c)
	log.Printf("%v\n", c.GetAttr(common.HPMax).GetIntValue())
	assert.Equal(t, 1, c.GetAttr(common.HPMax).GetIntValue())
}
