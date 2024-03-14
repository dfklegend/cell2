package builder

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"

	as "github.com/dfklegend/cell2/actorex/service"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/config"
	actormodule "github.com/dfklegend/cell2/node/modules/actor"
	"github.com/dfklegend/cell2/node/service"
	l "github.com/dfklegend/cell2/utils/logger"
)

// 方便service创建

type Builder struct {
	name     string
	producer func() actor.Actor
	postFunc func(*Builder, service.INodeServiceOwner)

	// temp
	cfg *config.ServiceInfo
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) WithName(name string) *Builder {
	b.name = name
	return b
}

func (b *Builder) GetName() string {
	return b.name
}

func (b *Builder) GetCfg() *config.ServiceInfo {
	return b.cfg
}

func (b *Builder) WithProducer(producer func() actor.Actor) *Builder {
	b.producer = producer
	return b
}

func (b *Builder) WithPostFunc(postFunc func(*Builder, service.INodeServiceOwner)) *Builder {
	b.postFunc = postFunc
	return b
}

func (b *Builder) BuildAndStart() {
	system := actormodule.GetSystem()
	rootContext := system.Root

	b.cfg = app.Node.GetServiceCfg(b.name)
	if b.cfg == nil {
		l.L.Errorf("can not find service cfg: %v", b.name)
		return
	}

	apiCollection := fmt.Sprintf("%v.remote", b.cfg.Type)
	props, ext := service.NewServiceWithDispatcher(b.producer,
		b.name, apiCollection)

	ext.WithPostFunc(func(s as.IService) {
		ns, _ := s.(service.INodeServiceOwner)
		b.postFunc(b, ns)
	})

	// restart的正确性
	ext.WithPostStartFunc(func(ctx actor.Context) {
		service.StartNodeService(rootContext, ctx.Self(), b.name, app.Node.GetServiceCfg(b.name))
	})
	rootContext.SpawnNamed(props, b.name)
}

/*
func Create(name string) {
	system := actormodule.GetSystem()
	rootContext := system.Root

	props, ext := service.NewServiceWithDispatcher(func() actor.Actor { return NewService() },
		name, "chat.remote")
	ext.WithPostFunc(func(s as.IService) {
		ns, _ := s.(service.INodeServiceOwner)
		impls.ServiceCreateCommonComponents(ns.GetNodeService(), app.Node.GetServiceCfg(name))
	})

	pid, _ := rootContext.SpawnNamed(props, name)
	service.StartNodeService(rootContext, pid, name, app.Node.GetServiceCfg(name))
}*/
