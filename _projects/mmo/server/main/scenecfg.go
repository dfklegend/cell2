package main

import (
	"mmo/modules/csv"
	"mmo/modules/csv/entry"
	"mmo/modules/scenecfg"
)

func LoadAllSceneCfgs(path string) {
	mgr := scenecfg.InitMgr(path)
	csv.Scene.Visit(func(item *entry.Scene) {
		mgr.Load(item.Id)
	})
}
