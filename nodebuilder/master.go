package builder

import (
	"github.com/dfklegend/cell2/apimapper/registry"
	"github.com/dfklegend/cell2/node/app"
	actormodule "github.com/dfklegend/cell2/node/modules/actor"
	clustermodule "github.com/dfklegend/cell2/node/modules/cluster"
	welcomemodule "github.com/dfklegend/cell2/node/modules/welcome"
	utils "github.com/dfklegend/cell2/nodeutils"
)

// 链式调用定义node的master初始化
// 便于使用者理解

type MasterBuilder struct {
}

func NewMasterBuilder() *MasterBuilder {
	return &MasterBuilder{}
}

func (b *MasterBuilder) ConfigLog(prefix, logDir string) *MasterBuilder {
	app.Node.EnableFileLog(prefix, logDir)
	return b
}

// RegisterAPIs 注册一下api集合
func (b *MasterBuilder) RegisterAPIs(needDump bool, cbDoRegister func()) *MasterBuilder {
	utils.NodeInitSystemAPI()
	if cbDoRegister != nil {
		cbDoRegister()
	}
	registry.Registry.Build()

	if needDump {
		registry.Registry.DumpAll()
	}
	return b
}

// Next 可以插入中间写一些额外的初始化代码
func (b *MasterBuilder) Next(cb func()) *MasterBuilder {
	cb()
	return b
}

func (b *MasterBuilder) StartMaster(configDir string, cbSucc func(succ bool)) *MasterBuilder {
	n := app.Node

	n.AddModule(welcomemodule.NewWelcomeModule())
	n.AddModule(actormodule.NewActorSystemModule())
	n.AddModule(clustermodule.NewClusterModule())

	n.Prepare(configDir)
	n.StartMaster(cbSucc)
	return b
}
