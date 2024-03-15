package csv

import (
	"mmo/modules/csv/base"
	"mmo/modules/csv/entry"
)

// ----

type equipCfg struct {
	*base.DataMgr[*entry.Equip]
}

func NewEquipCfg() *equipCfg {
	return &equipCfg{
		DataMgr: base.NewDataMgr[*entry.Equip](),
	}
}
