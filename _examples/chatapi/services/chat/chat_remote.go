package chat

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
	nodeapi.Registry.AddCollection("chat.remote").
		Register(&Entry{}, apientry.WithGroupName("chat"), apientry.WithNameFunc(strings.ToLower))
}

func Visit() {
}

type Entry struct {
	api.APIEntry
}

func (e *Entry) CreateRoom(ctx *as.RemoteContext, msg *mymsg.CSCreateRoom, cbFunc apientry.HandlerCBFunc) error {
	s := ctx.ActorContext.Actor().(*Service)
	log.Printf("createRoom in %v\n", s.Name)
	s.doCreateRoom(msg, cbFunc)
	return nil
}

func (e *Entry) Join(ctx *as.RemoteContext, msg *mymsg.ReqJoin, cbFunc apientry.HandlerCBFunc) error {
	s := ctx.ActorContext.Actor().(*Service)
	log.Printf("join in %v\n", s.Name)
	s.doJoin(msg, cbFunc)
	return nil
}

func (e *Entry) Say(ctx *as.RemoteContext, msg *mymsg.ReqChat, cbFunc apientry.HandlerCBFunc) error {
	s := ctx.ActorContext.Actor().(*Service)
	log.Printf("say in %v\n", s.Name)
	s.doChat(msg, cbFunc)
	return nil
}
