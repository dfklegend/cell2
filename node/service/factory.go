package service

import (
	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/node/config"
	l "github.com/dfklegend/cell2/utils/logger"
)

var (
	Factory = NewFactory()
)

//	Service构建器
type factory struct {
	// typeName: creator
	creators map[string]IServiceCreator
}

func NewFactory() *factory {
	return &factory{
		creators: make(map[string]IServiceCreator),
	}
}

func GetFactory() *factory {
	return Factory
}

func (f *factory) Register(typeName string, creator IServiceCreator) {
	f.creators[typeName] = creator
}

func (f *factory) Create(typeName string, name string) {
	creator := f.creators[typeName]
	if creator == nil {
		l.Log.Errorf("create %v failed, can not find service type:%v", name, typeName)
		return
	}
	creator.Create(name)
}

// StartNodeService "启动"nodeService
// 发送消息，通知启动(在service的routine里)
func StartNodeService(ctx *actor.RootContext, pid *actor.PID, name string,
	info *config.ServiceInfo, args ...interface{}) {
	msg := &StartServiceCmd{
		Name: name,
		Info: info,
		Args: args,
	}
	ctx.Send(pid, msg)
}
