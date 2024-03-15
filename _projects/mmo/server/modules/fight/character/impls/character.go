package charimpls

import (
	"github.com/dfklegend/cell2/utils/event/light"
	l "github.com/dfklegend/cell2/utils/logger"

	"mmo/common/entity"
	"mmo/modules/fight/attr"
	"mmo/modules/fight/attr/impls"
	"mmo/modules/fight/common"
)

type Character struct {
	common.IBufCtrl
	common.ISpecialStatusCtrl

	e           entity.IEntity
	id          common.CharId
	world       common.IWorld
	equipGroups []common.IEquipGroup

	attrs  *impls.AttrSlice
	events *light.EventCenter

	skill       common.ISkillCtrl
	skillTable  common.ISkillTable
	charWatcher common.ICharWatcher
	slots       common.IEquipSlots
	idTar       common.CharId

	normalAttackInterval int32
	attrCWatchers        []common.IAttrChangeWatcher

	// Lua proxy
	proxy *Proxy
}

func NewCharacter() *Character {
	c := &Character{
		equipGroups:          make([]common.IEquipGroup, 0),
		normalAttackInterval: 5000,
		attrCWatchers:        make([]common.IAttrChangeWatcher, 0),
		idTar:                -1,
	}
	c.ISpecialStatusCtrl = NewSpecialStatusCtrl(c)
	return c
}

func (c *Character) Init(id common.CharId, world common.IWorld, events *light.EventCenter) {
	c.id = id
	c.world = world
	c.events = events
	c.skillTable.Init(c, world.GetTimeProvider())
	c.proxy = newProxy(c)
}

func (c *Character) SetCharWatcher(watcher common.ICharWatcher) {
	c.charWatcher = watcher
}

func (c *Character) setBufCtrl(ctrl common.IBufCtrl) {
	c.IBufCtrl = ctrl
}

func (c *Character) GetProxy() common.ICharProxy {
	return c.proxy
}

func (c *Character) Start() {
	c.InitAttrWatchers()
}

func (c *Character) Destroy() {
	if c.IBufCtrl != nil {
		c.IBufCtrl.Destroy()
	}
}

func (c *Character) Update() {
	if c.skill == nil {
		return
	}
	c.skill.Update()
	if c.IBufCtrl != nil {
		c.IBufCtrl.Update()
	}
}

func (c *Character) GetId() common.CharId {
	return c.id
}

func (c *Character) GetChar(id common.CharId) common.ICharacter {
	return c.world.GetChar(id)
}

func (c *Character) GetPos() common.Pos {
	return common.Pos{}
}

func (c *Character) GetEvents() *light.EventCenter {
	return c.events
}

func (c *Character) GetWorld() common.IWorld {
	return c.world
}

func (c *Character) GetSkillTable() common.ISkillTable {
	return c.skillTable
}

func (c *Character) newAttrs(size int) {
	c.attrs = impls.NewAttrSlice(size)
}

func (c *Character) GetAttr(index int) attr.IAttr {
	return c.attrs.GetAttr(index)
}

func (c *Character) findEquipGroupIndex(group common.IEquipGroup) int {
	find := -1
	for i := 0; i < len(c.equipGroups); i++ {
		if c.equipGroups[i] == group {
			find = i
			break
		}
	}
	return find
}

func (c *Character) AddEquipGroup(group common.IEquipGroup) {
	if c.findEquipGroupIndex(group) != -1 {
		l.L.Errorf("equipgroup add already!!!")
		return
	}
	c.equipGroups = append(c.equipGroups, group)
	group.OnAdded(c)
	group.Equip(c)
}

func (c *Character) SetEquip(index int, id common.EquipId) {
	c.slots.SetEquip(index, id)
}

func (c *Character) RemoveEquip(index int) {
	c.slots.RemoveEquip(index)
}
