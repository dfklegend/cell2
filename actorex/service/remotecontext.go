package service

import "github.com/asynkron/protoactor-go/actor"

type RemoteContext struct {
	// 服务接口
	//Srv          *Service
	ActorContext actor.Context
}

func NewRemoteContext() *RemoteContext {
	return &RemoteContext{}
}

func (r *RemoteContext) Update(ctx actor.Context) {
	r.ActorContext = ctx
}

func (r *RemoteContext) Reserve() {
}

func (r *RemoteContext) Handle() {
}
