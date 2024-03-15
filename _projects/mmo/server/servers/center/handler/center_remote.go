package handler

import (
	"strings"

	as "github.com/dfklegend/cell2/actorex/service"
	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/apimapper/registry"
	l "github.com/dfklegend/cell2/utils/logger"

	"mmo/common/define"
	mymsg "mmo/messages"
)

func init() {
	registry.Registry.AddCollection("center.remote").
		Register(&Entry{}, apientry.WithName("centerremote"), apientry.WithNameFunc(strings.ToLower))
}

func Visit() {
}

type Entry struct {
	api.APIEntry
}

// RegisterCmdHandler(kickPrev)
// 上下线锁，防止上下线过程中重入()

func (e *Entry) ReqLogin(ctx *as.RemoteContext, msg *mymsg.CenterReqLogin, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	l.Log.Infof("ReqLogin in %v", s.Name)
	l.Log.Infof("msg: %+v", msg)

	s.Mgr.ReqLogin(msg.UId, msg.ServerId, msg.NetId, msg.KickPrev, cbFunc)

	//apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.CenterReqLoginAck{
	//	Code:        int32(code),
	//	IsReconnect: isReconnect,
	//})
}

func (e *Entry) OnSessionClose(ctx *as.RemoteContext, msg *mymsg.CenterOnSessionClose, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	l.Log.Infof("OnSessionClose in %v", s.Name)
	l.Log.Infof("msg: %+v", msg)

	s.Mgr.OnClientSessionClosed(msg.UId)

	apientry.CheckInvokeCBFunc(cbFunc, nil, nil)
}

func (e *Entry) OnLogicLogined(ctx *as.RemoteContext, msg *mymsg.OnLogicLogined, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	l.Log.Infof("OnLogicOnline in %v", s.Name)

	s.Mgr.OnLogicLogined(msg.UId, msg.LogicId)

	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
		Code: int32(define.Succ),
	})
}

func (e *Entry) OnLogicReonline(ctx *as.RemoteContext, msg *mymsg.OnLogicReOnline, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	l.Log.Infof("OnLogicReonline in %v", s.Name)

	s.Mgr.OnLogicReOnline(msg.UId)

	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
		Code: int32(define.Succ),
	})
}

func (e *Entry) ReqLogout(ctx *as.RemoteContext, msg *mymsg.ReqLogout, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	l.Log.Infof("ReqLogout in %v", s.Name)

	code := define.Succ
	if !s.Mgr.ReqLogout(msg.UId) {
		code = define.ErrFaild
	}
	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
		Code: int32(code),
	})
}

func (e *Entry) OnLogout(ctx *as.RemoteContext, msg *mymsg.OnLogout, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	s.GetLogger().Infof("%v OnLogicLogout ", msg.UId)

	s.Mgr.OnLogicLogout(msg.UId)

	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
		Code: int32(define.Succ),
	})
}

func (e *Entry) OnAbnormalLogout(ctx *as.RemoteContext, msg *mymsg.OnLogout, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	s.GetLogger().Infof("%v OnAbnormalLogout ", msg.UId)

	s.Mgr.OnLogicAbnormalLogout(msg.UId)

	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
		Code: int32(define.Succ),
	})
}

func (e *Entry) ReqSwitchLine(ctx *as.RemoteContext, msg *mymsg.ReqSwitchLine, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	s.GetLogger().Infof("ReqSwitchLine %v", msg.UId)

	code := define.Succ
	if !s.Mgr.ReqSwitchLine(msg.UId) {
		code = define.ErrFaild
	}

	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
		Code: int32(code),
	})
}

func (e *Entry) OnSwitchLineEnd(ctx *as.RemoteContext, msg *mymsg.OnSwitchLineEnd, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	s.GetLogger().Infof("OnSwitchLineEnd %v %v", msg.UId, msg.Succ)

	code := define.Succ
	if !s.Mgr.OnSwitchLineEnd(msg.UId, msg.Succ) {
		code = define.ErrFaild
	}

	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
		Code: int32(code),
	})
}
