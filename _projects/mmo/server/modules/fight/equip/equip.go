package equip

import (
	"mmo/modules/csv/entry"
	"mmo/modules/csv/structs"
	"mmo/modules/fight/common"
)

// Item 装备物品
type Item struct {
	id  common.EquipId
	cfg *entry.Equip
}

func NewItem(cfg *entry.Equip) *Item {
	return &Item{
		id:  cfg.Id,
		cfg: cfg,
	}
}

func (i *Item) Equip(tar common.ICharacter) {
	i.offsetAttrs(tar, true)
	i.addBuf(tar)
}

func (i *Item) Unequip(tar common.ICharacter) {
	i.offsetAttrs(tar, false)
	i.removeBuf(tar)
}

func (i *Item) offsetAttrs(tar common.ICharacter, add bool) {
	cfg := i.cfg

	if cfg.Attr0.Type != -1 {
		i.offsetAttr(tar, add, cfg.Attr0, cfg.V0)
	}
	if cfg.Attr1.Type != -1 {
		i.offsetAttr(tar, add, cfg.Attr1, cfg.V1)
	}
	if cfg.Attr2.Type != -1 {
		i.offsetAttr(tar, add, cfg.Attr2, cfg.V2)
	}
	if cfg.Attr3.Type != -1 {
		i.offsetAttr(tar, add, cfg.Attr3, cfg.V3)
	}
}

// TODO: 考虑多个装备同一个buf的情况

func (i *Item) offsetAttr(tar common.ICharacter, add bool, attr structs.AttrValue, v float32) {
	if attr.Type == -1 {
		return
	}

	final := v

	if !add {
		final = -final
	}

	if attr.IsPercent {
		tar.OffsetPercent(attr.Type, final)
		return
	}
	tar.OffsetBase(attr.Type, float64(final))
}

func (i *Item) addBuf(tar common.ICharacter) {
	cfg := i.cfg
	if cfg.Buf == "" {
		return
	}
	tar.AddBuf(tar, tar, cfg.Buf, 1, 1)
}

func (i *Item) removeBuf(tar common.ICharacter) {
	cfg := i.cfg
	if cfg.Buf == "" {
		return
	}
	tar.RemoveBuf(cfg.Buf)
}
