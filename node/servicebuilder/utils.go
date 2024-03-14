package builder

import (
	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/node/client/impls"
	"github.com/dfklegend/cell2/node/client/impls/pomelo"
	"github.com/dfklegend/cell2/node/service"
)

// StartFrontService 启动标准前端服务
// 映射到api集 serviceType.remote
func StartFrontService(name string, producer func() actor.Actor) {
	NewBuilder().
		WithName(name).
		WithProducer(producer).
		WithPostFunc(func(b *Builder, nsOwner service.INodeServiceOwner) {
			impls.ServiceCreateCommonComponents(nsOwner.GetNodeService(), b.GetCfg())
			pomelo.ServiceCreateAcceptors(nsOwner.GetNodeService(), name, b.GetCfg())
		}).
		BuildAndStart()
}

// StartBackService 启动标准后端服务
func StartBackService(name string, producer func() actor.Actor) {
	NewBuilder().
		WithName(name).
		WithProducer(producer).
		WithPostFunc(func(b *Builder, nsOwner service.INodeServiceOwner) {
			impls.ServiceCreateCommonComponents(nsOwner.GetNodeService(), b.GetCfg())
		}).
		BuildAndStart()
}
