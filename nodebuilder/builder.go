package builder

import (
	"github.com/dfklegend/cell2/apimapper/registry"
	"github.com/dfklegend/cell2/baseapp"
	"github.com/dfklegend/cell2/baseapp/interfaces"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/route"
	"github.com/dfklegend/cell2/node/service"
	utils "github.com/dfklegend/cell2/nodeutils"
)

// 链式调用定义node的初始化顺序
// 便于使用者理解

type Builder struct {
	onPrepare func()
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) ConfigLog(prefix, logDir string) *Builder {
	app.Node.EnableFileLog(prefix, logDir)
	return b
}

// SetDefaultLaunchFunc 设置缺省启动函数
func (b *Builder) SetDefaultLaunchFunc(f func(app interfaces.IApp)) *Builder {
	baseapp.SetDefaultLaunchFunc(f)
	return b
}

// RegisterLaunchFunc 注册一个启动模式
func (b *Builder) RegisterLaunchFunc(name string, f func(app interfaces.IApp)) *Builder {
	baseapp.RegisterLaunchFunc(name, f)
	return b
}

// RegisterLaunchModes 注册其他的启动模式
func (b *Builder) RegisterLaunchModes(cb func()) *Builder {
	cb()
	return b
}

/*
// 注册launch mode
	func() {
		baseapp.LaunchFunc("allinone", func(app interfaces.IApp) {
			log.Printf("allinone mode\n")
			utils.NodeAddCommonModules(app)
			utils.NodeAddClusterModules(app)
		})
	}
*/

// RegisterAPIs 注册一下api集合
func (b *Builder) RegisterAPIs(needDump bool, cbDoRegister func()) *Builder {
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

// RegisterServiceCreators 注册一下服务构建器
func (b *Builder) RegisterServiceCreators(cb func(factory service.IServiceFactory)) *Builder {
	cb(service.GetFactory())
	// factory.Register("db", db.NewCreator())
	return b
}

// RegisterRoutes 注册一下服务路由函数
func (b *Builder) RegisterRoutes(cb func(rs *route.RouteService)) *Builder {
	cb(route.GetRouteService())
	return b
}

// Next 可以插入中间写一些额外的初始化代码
func (b *Builder) Next(cb func()) *Builder {
	cb()
	return b
}

func (b *Builder) OnPrepare(cb func()) *Builder {
	b.onPrepare = cb
	return b
}

func (b *Builder) StartApp(configDir string, nodeId string, cbSucc func(succ bool)) *Builder {
	n := app.Node
	n.Prepare(configDir)
	if b.onPrepare != nil {
		b.onPrepare()
	}
	n.StartNode(nodeId, cbSucc)
	return b
}
