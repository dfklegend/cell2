package impls

import (
	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/builtin/msgs"
	"github.com/dfklegend/cell2/node/client"
	"github.com/dfklegend/cell2/node/client/impls/config"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/logger"
)

// 	serviceType.category.method
func SplitClientRoute(route string) (string, string, string) {
	return app.SplitClientRoute(route)
}

func PushMessageById(ns *service.NodeService, serverId string, sessionId uint32, route string, msg any) {
	pid := app.GetServicePID(serverId)
	if pid == nil {
		logger.Log.Errorf("can not find service: %v", serverId)
		return
	}

	ids := make([]uint32, 1)
	ids[0] = sessionId
	pushMessageByIds(ns, pid, ids, route, msg)
}

func PushMessageByIds(ns *service.NodeService, serverId string, ids []uint32, route string, msg any) {
	pid := app.GetServicePID(serverId)
	if pid == nil {
		logger.Log.Errorf("can not find service: %v", serverId)
		return
	}
	pushMessageByIds(ns, pid, ids, route, msg)
}

func pushMessageByIds(ns *service.NodeService, pid *actor.PID, ids []uint32, route string, msg any) {
	pmsg := &msgs.PushMsg{}
	pmsg.Route = route
	pmsg.Ids = ids
	pmsg.Data, _ = config.GetConfig().Serializer.Marshal(msg)
	ns.RequestEx(pid, "sys.pushmsg", pmsg, nil)
}

func AddOnSessionOnClose(ns *service.NodeService, netId uint32, onClose client.CBSessionOnClose) {
	ns.GetComponent("handler").(*HandlerComponent).AddOnSessionClose(netId, onClose)
}
