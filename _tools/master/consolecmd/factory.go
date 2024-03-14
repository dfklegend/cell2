package consolecmd

import (
	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/actorex/service"
)

func CreateService(system *actor.ActorSystem) *actor.PID {
	rootContext := system.Root

	name := "consolecmd"
	props, ext := service.NewServicePropsWithNewScheDisp(func() actor.Actor { return NewService() },
		name)
	ext.WithAPIs("consolecmd")
	pid, _ := rootContext.SpawnNamed(props, name)
	return pid
}
