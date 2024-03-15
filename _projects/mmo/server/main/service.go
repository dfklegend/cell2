package main

import (
	"github.com/dfklegend/cell2/node/service"

	"mmo/servers/center/handler"
	dbhandler "mmo/servers/db/handler"
	"mmo/servers/gate/handler"
	ls "mmo/servers/logic/logicservice"
	"mmo/servers/logic_old/handler"
	sceneservice "mmo/servers/scene/handler"
	"mmo/servers/scene/logic/logics"
	"mmo/servers/scenem/handler"
)

func registerServices(factory service.IServiceFactory) {
	// 注册服务构建器
	factory.Register("center", handler.NewCreator())
	factory.Register("db", dbhandler.NewCreator())
	factory.Register("gate", gate.NewCreator())
	factory.Register("logic_old", logic.NewCreator())
	factory.Register("logic", ls.NewCreator())
	factory.Register("scenem", scenem.NewCreator())
	factory.Register("scene", sceneservice.NewCreator())
}

func registerServiceBonus() {
	scenelogics.Visit()
}
