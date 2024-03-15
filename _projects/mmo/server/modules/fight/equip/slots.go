package equip

import (
	l "github.com/dfklegend/cell2/utils/logger"

	"mmo/modules/csv"
	"mmo/modules/fight/common"
)

// Slots 可以装备多个装备
type Slots struct {
	owner common.ICharacter
	items []*Item
}

func NewSlots(num int) *Slots {
	return &Slots{
		items: make([]*Item, num),
	}
}

func (s *Slots) OnAdded(owner common.ICharacter) {
	s.owner = owner
}

func (s *Slots) Equip(owner common.ICharacter) {
	for _, v := range s.items {
		if v != nil {
			v.Equip(owner)
		}
	}
}

func (s *Slots) Unequip(owner common.ICharacter) {
	for _, v := range s.items {
		if v != nil {
			v.Unequip(owner)
		}
	}
}

func (s *Slots) isValidSlot(index int) bool {
	if index < 0 || index >= len(s.items) {
		return false
	}
	return true
}

func (s *Slots) SetEquip(index int, id common.EquipId) {
	if !s.isValidSlot(index) {
		return
	}

	cfg := csv.Equip.GetEntry(id)
	if cfg == nil {
		l.L.Warnf("error equip %v at slot %v", id, index)
		return
	}

	old := s.items[index]
	if old != nil {
		old.Unequip(s.owner)
	}

	item := NewItem(cfg)
	item.Equip(s.owner)
	s.items[index] = item
}

func (s *Slots) RemoveEquip(index int) {
	if !s.isValidSlot(index) {
		return
	}

	old := s.items[index]
	if old != nil {
		old.Unequip(s.owner)
		s.items[index] = nil
	}
}
