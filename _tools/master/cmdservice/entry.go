package cmdservice

import (
	"github.com/dfklegend/cell2/actorex/service"
	"github.com/dfklegend/cell2/actorex/service/servicemsgs"
	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/apimapper/registry"
	"master/messages"
)

func RegisterEntry() {
	registry.RegisterWithLowercaseName(&Entry{}, "cmdservice", "service")
}

type Entry struct {
	api.APIEntry
}

func (e *Entry) Start(ctx *service.RemoteContext, msg *servicemsgs.EmptyArg) {
	s := ctx.ActorContext.Actor().(*Service)
	s.Start()
}

func (e *Entry) Cmd(ctx *service.RemoteContext, msg *messages.Cmd, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)

	s.ProcessCmd(msg.Cmd, func(result string) {
		apientry.CheckInvokeCBFunc(cbFunc, nil, &messages.CmdAck{
			Result: result,
		})
	})

}
