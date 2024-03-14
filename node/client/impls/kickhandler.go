package impls

import (
	"github.com/dfklegend/cell2/node/service"
)

type IKickHandler interface {
	HandleKick(ns *service.NodeService, sessions *ClientSessions, netId uint32)
}
