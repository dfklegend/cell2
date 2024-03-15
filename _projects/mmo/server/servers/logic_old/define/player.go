package define

import (
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/event/light"

	"mmo/servers/logic_old/systems/bridge"
)

// ILogicPlayer
// 	fightmode
// 	system
type ILogicPlayer interface {
	GetUId() int64

	SetDirt()
	GetNodeService() *service.NodeService
	GetEvents() *light.EventCenter

	// systems

	GetCharCard() bridge.ICharCard

	PushMsg(route string, msg interface{})
}
