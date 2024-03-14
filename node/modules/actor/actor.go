package actormodule

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"

	"github.com/dfklegend/cell2/baseapp/interfaces"
	"github.com/dfklegend/cell2/baseapp/module"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/utils/common"
	l "github.com/dfklegend/cell2/utils/logger"
)

var (
	system *actor.ActorSystem
)

func GetSystem() *actor.ActorSystem {
	return system
}

type ActorSystemModule struct {
	*module.BaseModule
}

func NewActorSystemModule() *ActorSystemModule {
	return &ActorSystemModule{
		module.NewBaseModule(),
	}
}

func (a *ActorSystemModule) Start(next interfaces.FuncWithSucc) {
	app := app.Node
	info := app.GetNodeInfo()
	if info == nil {
		l.Log.Errorf("actor system module can not find node info!")
		next(false)
		return
	}

	l.Log.Infof("init actorSystem: %v", info.Address)
	ip, port := common.SplitAddress(info.Address)

	system = actor.NewActorSystem()
	config := remote.Configure(ip, port)
	remoter := remote.NewRemote(system, config)
	remoter.Start()

	app.SetActorSystem(system)
	next(true)
}

func (a *ActorSystemModule) Stop(next interfaces.FuncWithSucc) {
	system.Shutdown()
	next(true)
}
