package channel

import (
	"github.com/dfklegend/cell2/node/service"
)

// IPushMessager 避免交叉引用
var (
	pushImpl IPushMessager
)

func SetPushImpl(i IPushMessager) {
	pushImpl = i
}

func GetPushImpl() IPushMessager {
	return pushImpl
}

//	IPushMessager 解除依赖引用
type IPushMessager interface {
	PushMessageById(ns *service.NodeService, serverId string, sessionId uint32, route string, msg any)
	PushMessageByIds(ns *service.NodeService, serverId string, ids []uint32, route string, msg any)
}
