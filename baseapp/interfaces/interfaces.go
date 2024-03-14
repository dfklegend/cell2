package interfaces

import "github.com/dfklegend/cell2/utils/runservice"

type FuncWithSucc func(succ bool)

//	app模块，通过module来组织实际app的启动流程
//	根据添加顺序依次start，反向stop
type IAppModule interface {
	Init(rs *runservice.StandardRunService)
	Start(next FuncWithSucc)
	Stop(next FuncWithSucc)
}

//	App
//	拥有一个独立的携程
//	主要负责App的启动与关闭
type IApp interface {
	AddModule(module IAppModule)
	Start(finish FuncWithSucc)
	Stop(finish FuncWithSucc)
}

/*
	IBlackBoard
	黑板，记录一些共享变量
*/
type IBlackBoard interface {
	GetValue(key string, def interface{}) interface{}
	SetValue(string, interface{})
}

//	启动模式
//	通过实现不同的启动模式，来规划节点功能
// 	app.AddModule
//	app.Start
type ILaunchMode interface {
	PrepareModules(app IApp)
}
