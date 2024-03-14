package impls

import (
	"github.com/dfklegend/cell2/node/builtin/channel"
	"github.com/dfklegend/cell2/node/service"
)

func init() {
	channel.SetPushImpl(&ChannelPushMessageImpl{})
}

type ChannelPushMessageImpl struct {
}

func (p *ChannelPushMessageImpl) PushMessageById(ns *service.NodeService, serverId string, sessionId uint32, route string, msg any) {
	PushMessageById(ns, serverId, sessionId, route, msg)
}

func (p *ChannelPushMessageImpl) PushMessageByIds(ns *service.NodeService, serverId string, ids []uint32, route string, msg any) {
	PushMessageByIds(ns, serverId, ids, route, msg)
}
