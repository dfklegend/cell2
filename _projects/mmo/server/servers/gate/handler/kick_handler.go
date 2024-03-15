package gate

import (
	"time"

	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/client/impls"
	"github.com/dfklegend/cell2/node/service"

	"mmo/messages/cproto"
)

type kickHandler struct {
	impls.IKickHandler
}

func newKickHandler() *kickHandler {
	return &kickHandler{}
}

func (h *kickHandler) HandleKick(ns *service.NodeService, sessions *impls.ClientSessions, netId uint32) {
	ns.GetLogger().Infof("kickHandler HandleKick: %v", netId)

	// 发个消息
	app.PushMessageById(ns, ns.Name, netId, "kick", &cproto.Kick{
		Reason: "kick",
	})
	// 等一秒再关
	ns.GetRunService().GetTimerMgr().After(time.Second, func(args ...any) {
		sessions.DoKick(netId)
	})
}
