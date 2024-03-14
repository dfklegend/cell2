package service

import (
	"github.com/asynkron/protoactor-go/actor"

	as "github.com/dfklegend/cell2/actorex/service"
	"github.com/dfklegend/cell2/apimapper/registry"
	"github.com/dfklegend/cell2/utils/logger"
)

// NewServiceWithDispatcher
// 定义remote的接口处理者
// handler的接口由handlerComponent创建时根据服务类型来指定(并不是服务自身的接口，也不能直接调用)
// apiCollection  需要注册的collectioncollection
func NewServiceWithDispatcher(producer actor.Producer, name, apiCollection string) (*actor.Props, *as.ExtProps) {
	props, ext := as.NewServicePropsWithNewScheDisp(producer, name)
	ext.WithDispatcher(NewDispatcher(apiCollection))
	return props, ext
}

func NewDispatcher(collections string) *as.APIDispatcher {
	d := as.NewDispatcher(registry.Registry.GetCollection(collections),
		registry.Registry.GetCollection(SystemAPI))
	return d
}

func EnableServiceLogFormatter() {
	logger.GetLogProxy("default").SetFormatter(NewServiceLogFormat("default"))
}
