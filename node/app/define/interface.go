package define

import (
	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/node/config"
)

type INodeApp interface {
	GetActorSystem() *actor.ActorSystem

	GetService(name string) *actor.PID
	FilterSelfServices(filter func(name string, cfg *config.ServiceInfo))

	UpdateNodeState(state int)
	StopNode(fin func(succ bool))
}
