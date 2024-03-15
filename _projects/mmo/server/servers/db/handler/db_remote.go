package handler

import (
	"strings"

	as "github.com/dfklegend/cell2/actorex/service"
	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/apimapper/registry"
	l "github.com/dfklegend/cell2/utils/logger"

	mymsg "mmo/messages"
)

func init() {
	registry.Registry.AddCollection("db.remote").
		Register(&Entry{}, apientry.WithName("dbremote"), apientry.WithNameFunc(strings.ToLower))
}

func Visit() {
}

type Entry struct {
	api.APIEntry
}

func (e *Entry) Auth(ctx *as.RemoteContext, msg *mymsg.DBAuth, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	l.Log.Infof("entry in %v", s.Name)
	l.Log.Infof("msg: %+v", msg)

	redis := s.redis
	s.RequestEx(redis, "redis.auth", msg, func(err error, raw interface{}) {
		l.L.Infof("got auth ret: %v", raw)
		apientry.CheckInvokeCBFunc(cbFunc, nil, raw)
	})
}

func (e *Entry) LoadPlayer(ctx *as.RemoteContext, msg *mymsg.DBLoadPlayer, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	l.Log.Infof("entry in %v", s.Name)
	l.Log.Infof("msg: %+v", msg)

	redis := s.redis
	s.RequestEx(redis, "redis.loadplayer", msg, func(err error, raw interface{}) {
		l.L.Infof("got loadplayer ret: %v", raw)
		apientry.CheckInvokeCBFunc(cbFunc, err, raw)
	})
}

func (e *Entry) SavePlayer(ctx *as.RemoteContext, msg *mymsg.DBSavePlayer, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	l.Log.Infof("SavePlayer in %v", s.Name)
	l.Log.Infof("msg: %+v", msg)

	redis := s.redis
	s.RequestEx(redis, "redis.saveplayer", msg, func(err error, raw interface{}) {
		l.L.Infof("got saveplayer ret: %v", raw)
		apientry.CheckInvokeCBFunc(cbFunc, err, raw)
	})
}
