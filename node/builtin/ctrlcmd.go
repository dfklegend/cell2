package builtin

import (
	"strings"

	as "github.com/dfklegend/cell2/actorex/service"
	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/apimapper/registry"
	"github.com/dfklegend/cell2/node/builtin/msgs"
	"github.com/dfklegend/cell2/node/service"
	l "github.com/dfklegend/cell2/utils/logger"
)

func init() {
	registry.Registry.AddCollection(service.SystemAPI).
		Register(&CtrlEventEntry{}, apientry.WithGroupName("ctrl"), apientry.WithNameFunc(strings.ToLower))
}

// CtrlEventEntry
// 处理控制命令的接口集
type CtrlEventEntry struct {
	api.APIEntry
}

//	处理ctrl.cmd
func (e *CtrlEventEntry) Cmd(ctx *as.RemoteContext, msg *msgs.CtrlCmd, cbFunc apientry.HandlerCBFunc) {
	// 如果异常，那么就是用户service没有实现INodeServiceOwner
	s := ctx.ActorContext.Actor().(service.INodeServiceOwner)
	ns := s.GetNodeService()

	listener := ns.GetCtrlCmdListener()
	if listener == nil {
		l.Log.Infof("%v has no CtrlEventListener", ns.Name)

		apientry.CheckInvokeCBFunc(cbFunc, nil, &msgs.CtrlCmdAck{
			Result: "no listener",
		})
		return
	}

	apientry.CheckInvokeCBFunc(cbFunc, nil, &msgs.CtrlCmdAck{
		Result: listener.Handler(msg.Cmd),
	})
	return
}
