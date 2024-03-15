package handler

import (
	"strings"

	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/apimapper/registry"
	"github.com/dfklegend/cell2/node/client/impls"
	l "github.com/dfklegend/cell2/utils/logger"

	"mmo/messages/cproto"
	"mmo/servers/logic/logicservice"
	"mmo/utils"
)

func init() {
	registry.Registry.AddCollection("logic.handler").
		Register(&Handler{}, apientry.WithName("logic"), apientry.WithNameFunc(strings.ToLower))
}

type Handler struct {
	api.APIEntry
}

// StartFight 请求战斗
func (e *Handler) StartFight(ctx *impls.HandlerContext, msg *cproto.StartGame, cbFunc apientry.HandlerCBFunc) {
	// . 检查当前状态是否允许战斗
	// . 请求分配战斗服
	// . 请求进入战斗服
	s := ctx.ActorContext.Actor().(*logicservice.Service)
	l.L.Infof("StartFight in %v", s.Name)

	uid := utils.GetIdFromContext(ctx)
	if uid == 0 {
		l.L.Errorf("error find uid from ctx")
		return
	}
}

//func (e *Handler) Cmd(ctx *impls.HandlerContext, msg *cproto.Cmd, cbFunc apientry.HandlerCBFunc) {
//	s := ctx.ActorContext.Actor().(*logicservice.Service)
//	sid := ctx.GetSession().GetID()
//	uid := convert.TryParseInt64(sid, 0)
//	if uid == 0 {
//		l.L.Errorf("cmd error uid")
//		return
//	}
//	l.L.Infof("cmd: %v", s.Name)
//
//	//cmdCtx := s.cmdContext
//	//cmdCtx.UId = uid
//	//cmdCtx.P = s.Mgr.GetPlayer(uid)
//	//s.cmdMgr.Dispatch(cmdCtx, msg.Cmd, func(result string) {
//	//	apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.CmdAck{
//	//		Result: result,
//	//	})
//	//})
//}

func (e *Handler) SystemCmd(ctx *impls.HandlerContext, msg *cproto.ReqSystemCmd, cbFunc apientry.HandlerCBFunc) {

}
