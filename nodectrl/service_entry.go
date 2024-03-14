package nodectrl

import (
	"strings"

	"github.com/dfklegend/cell2/actorex/service"
	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/apimapper/registry"
	"github.com/dfklegend/cell2/node/builtin/msgs"
)

func init() {
	registry.Registry.AddCollection("node.admin").
		Register(&AdminEntry{}, apientry.WithGroupName("ctrl"), apientry.WithNameFunc(strings.ToLower))
}

type AdminEntry struct {
	api.APIEntry
}

// Cmd call by master， 控制节点
func (e *AdminEntry) Cmd(ctx *service.RemoteContext, msg *msgs.CtrlCmd, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*AdminService)
	apientry.CheckInvokeCBFunc(cbFunc, nil, &msgs.CtrlCmdAck{
		Result: s.ProcessCmd(msg.Cmd),
	})
}

// ServiceCmd from service
func (e *AdminEntry) ServiceCmd(ctx *service.RemoteContext, msg *msgs.ServiceCmd, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*AdminService)
	apientry.CheckInvokeCBFunc(cbFunc, nil, &msgs.ServiceCmdAck{
		Result: s.ProcessServiceCmd(msg.Name, msg.Cmd),
	})
}
