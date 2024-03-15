package main

import (
	"mmo/servers/center/handler"
	dbhandler "mmo/servers/db/handler"
	"mmo/servers/gate/handler"
	logichandler "mmo/servers/logic/handler"
	"mmo/servers/logic_old/handler"
	sceneservice "mmo/servers/scene/handler"
	"mmo/servers/scenem/handler"
)

// RegisterAllUserEntries
// 访问注册每个服务的接口集
func RegisterAllUserEntries() {
	gate.Visit()
	handler.Visit()
	dbhandler.Visit()
	logic.Visit()
	sceneservice.Visit()
	scenem.Visit()
	logichandler.Visit()
}

func RegisterAllAPI() {
	RegisterAllUserEntries()
}
