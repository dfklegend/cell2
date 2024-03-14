package channel

import (
	"github.com/dfklegend/cell2/node/service"
)

func PushMessageById(ns *service.NodeService, serverId string, sessionId uint32, route string, msg any) {
	GetPushImpl().PushMessageById(ns, serverId, sessionId, route, msg)
}

func PushMessageByIds(ns *service.NodeService, serverId string, ids []uint32, route string, msg any) {
	GetPushImpl().PushMessageByIds(ns, serverId, ids, route, msg)
}
