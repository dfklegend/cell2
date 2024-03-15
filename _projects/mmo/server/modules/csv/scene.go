package csv

import (
	"mmo/modules/csv/base"
	"mmo/modules/csv/entry"
)

type sceneCfg struct {
	*base.IntDataMgr[*entry.Scene]
}

func NewSceneCfg() *sceneCfg {
	return &sceneCfg{
		IntDataMgr: base.NewIntDataMgr[*entry.Scene](),
	}
}
