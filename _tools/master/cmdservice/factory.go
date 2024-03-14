package cmdservice

import (
	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/actorex/service"
)

func CreateService(system *actor.ActorSystem) *actor.PID {
	rootContext := system.Root

	name := "cmdservice"
	props, ext := service.NewServicePropsWithNewScheDisp(func() actor.Actor { return NewService() },
		name)
	ext.WithAPIs("cmdservice")
	pid, _ := rootContext.SpawnNamed(props, name)
	return pid
}
