package csv

import (
	"mmo/modules/csv/base"
	"mmo/modules/csv/entry"
)

// ----

type bufCfg struct {
	*base.DataMgr[*entry.Buf]
}

func NewBufCfg() *bufCfg {
	return &bufCfg{
		DataMgr: base.NewDataMgr[*entry.Buf](),
	}
}
