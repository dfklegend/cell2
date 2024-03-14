package chat

import (
	"strings"

	mymsg "chat2/messages"
	as "github.com/dfklegend/cell2/actorex/service"
	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	nodeapi "github.com/dfklegend/cell2/apimapper/registry"
	l "github.com/dfklegend/cell2/utils/logger"
)

func init() {
	nodeapi.Registry.AddCollection("chat.remote").
		Register(&Entry{}, apientry.WithGroupName("chatremote"), apientry.WithNameFunc(strings.ToLower))
}

func Visit() {
}

type Entry struct {
	api.APIEntry
}

func (e *Entry) Entry(ctx *as.RemoteContext, msg *mymsg.RoomEntry, cbFunc apientry.HandlerCBFunc) error {
	s := ctx.ActorContext.Actor().(*Service)
	l.Log.Infof("entry in %v\n", s.Name)
	l.Log.Infof("msg: %+v\n", msg)

	s.cs.PlayerEnter(msg.UId, msg.Name, msg.ServerId, msg.NetId)
	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{})
	return nil
}

func (e *Entry) Leave(ctx *as.RemoteContext, msg *mymsg.RoomLeave, cbFunc apientry.HandlerCBFunc) error {
	s := ctx.ActorContext.Actor().(*Service)
	l.Log.Infof("leave in %v\n", s.Name)
	l.Log.Infof("msg: %+v\n", msg)

	s.cs.PlayerLeave(msg.UId)

	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{})
	return nil
}
