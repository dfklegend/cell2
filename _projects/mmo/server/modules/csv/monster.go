package csv

import (
	"mmo/modules/csv/base"
	"mmo/modules/csv/entry"
)

type monsterTemplateCfg struct {
	*base.DataMgr[*entry.MonsterTemplate]
}

func NewMonsterTemplateCfg() *monsterTemplateCfg {
	return &monsterTemplateCfg{
		DataMgr: base.NewDataMgr[*entry.MonsterTemplate](),
	}
}

type monsterCfg struct {
	*base.DataMgr[*entry.Monster]
}

func NewMonsterCfg() *monsterCfg {
	return &monsterCfg{
		DataMgr: base.NewDataMgr[*entry.Monster](),
	}
}
