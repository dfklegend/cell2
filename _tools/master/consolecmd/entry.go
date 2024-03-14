package consolecmd

import (
	"github.com/dfklegend/cell2/actorex/service"
	"github.com/dfklegend/cell2/actorex/service/servicemsgs"
	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/registry"
)

func RegisterEntry() {
	registry.RegisterWithLowercaseName(&Entry{}, "consolecmd", "console")
}

type Entry struct {
	api.APIEntry
}

func (e *Entry) Start(ctx *service.RemoteContext, msg *servicemsgs.EmptyArg) {
	s := ctx.ActorContext.Actor().(*Service)
	s.Start()
}
