package chatmgr

import (
	"log"
	"strings"

	mymsg "chatapi/messages"
	as "github.com/dfklegend/cell2/actorex/service"
	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	nodeapi "github.com/dfklegend/cell2/apimapper/registry"
)

func init() {
	nodeapi.Registry.AddCollection("chatm.remote").
		Register(&Entry{}, apientry.WithGroupName("chatm"), apientry.WithNameFunc(strings.ToLower))
}

func Visit() {
}

type Entry struct {
	api.APIEntry
}

func (e *Entry) RefreshInfo(ctx *as.RemoteContext, msg *mymsg.CSRefreshInfo) error {
	s := ctx.ActorContext.Actor().(*Service)
	//log.Printf("RefreshInfo in %v\n", s.Name)
	s.mgr.OnServiceInfo(msg.ChatServiceId, msg.RoomNum, msg.PlayerNum)
	return nil
}

func (e *Entry) RoomStat(ctx *as.RemoteContext, msg *mymsg.CSRoomStat) error {
	s := ctx.ActorContext.Actor().(*Service)
	log.Printf("RoomStat in %v\n", s.Name)
	s.mgr.updateRoomStat(msg.RoomID, msg.PlayerNum)
	return nil
}

func (e *Entry) ReqLogin(ctx *as.RemoteContext, msg *mymsg.CMReqLogin, cbFunc apientry.HandlerCBFunc) error {
	s := ctx.ActorContext.Actor().(*Service)
	log.Printf("ReqLogin in %v\n", s.Name)
	s.mgr.Login(s, msg, cbFunc)
	return nil
}
