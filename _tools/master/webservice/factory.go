package webservice

import (
	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/actorex/service"
)

func CreateService(system *actor.ActorSystem) *actor.PID {
	rootContext := system.Root

	name := "webservice"
	props, ext := service.NewServicePropsWithNewScheDisp(
		func() actor.Actor {
			return NewService()
		},
		name,
	)
	ext.WithAPIs("webservice")
	pid, _ := rootContext.SpawnNamed(props, name)
	return pid
}

func NewService() *WebService {
	s := &WebService{
		Service: service.NewService(),
	}
	s.Service.InitReqReceiver(s)
	return s
}
