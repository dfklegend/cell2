package context

import (
	"github.com/dfklegend/cell2/node/service"

	"mmo/servers/logic_old"
)

type CmdContext struct {
	NS  *service.NodeService
	UId int64
	P   *logic_old.Player
}
