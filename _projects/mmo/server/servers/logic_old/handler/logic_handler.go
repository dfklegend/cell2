package logic

import (
	"strings"

	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/apimapper/registry"
	"github.com/dfklegend/cell2/node/client/impls"
	"github.com/dfklegend/cell2/utils/convert"
	l "github.com/dfklegend/cell2/utils/logger"

	"mmo/common/define"
	"mmo/messages/cproto"
	"mmo/servers/logic_old/fightmode"
	"mmo/servers/logic_old/fightmode/builder"
	"mmo/utils"
)

func init() {
	registry.Registry.AddCollection("logic_old.handler").
		Register(&Handler{}, apientry.WithName("logic_old"), apientry.WithNameFunc(strings.ToLower))
}

type Handler struct {
	api.APIEntry
}

// StartFight 请求战斗
func (e *Handler) StartFight(ctx *impls.HandlerContext, msg *cproto.StartGame, cbFunc apientry.HandlerCBFunc) {
	// . 检查当前状态是否允许战斗
	// . 请求分配战斗服
	// . 请求进入战斗服
	s := ctx.ActorContext.Actor().(*Service)
	l.L.Infof("StartFight in %v", s.Name)

	uid := utils.GetIdFromContext(ctx)
	if uid == 0 {
		l.L.Errorf("error find uid from ctx")
		return
	}

	player := s.Mgr.GetPlayer(uid)
	if player == nil {
		l.L.Errorf("can not find player: %v", uid)
		apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.NormalAck{
			Code: int32(define.ErrLogicCannotFindPlayer),
			Err:  "failed",
		})
		return
	}

	if player.GetCharCardNum() < 2 {
		apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.NormalAck{
			Code: int32(define.ErrWithStr),
			Err:  "角色卡不足",
		})
		return
	}

	// create fight mode
	var mode fightmode.ISceneFightMode
	mode = builder.BuildCardFight(player, player.GetCard(msg.DownCard), player.GetCard(msg.UpCard))
	player.SetCurFight(mode)
	player.StartFight(mode, s.GetNodeService(), func(succ bool) {
		if succ {
			apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.NormalAck{
				Code: int32(define.Succ),
			})
		} else {
			player.SetCurFight(nil)
			apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.NormalAck{
				Code: int32(define.ErrWithStr),
				Err:  "failed",
			})
		}
	})
}

// Cmd
// 提供了简单的字符串命令协议
func (e *Handler) Cmd(ctx *impls.HandlerContext, msg *cproto.Cmd, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	sid := ctx.GetSession().GetID()
	uid := convert.TryParseInt64(sid, 0)
	if uid == 0 {
		l.L.Errorf("cmd error uid")
		return
	}

	cmdCtx := s.cmdContext
	cmdCtx.UId = uid
	cmdCtx.P = s.Mgr.GetPlayer(uid)
	s.cmdMgr.Dispatch(cmdCtx, msg.Cmd, func(result string) {
		apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.CmdAck{
			Result: result,
		})
	})
}

func (e *Handler) SystemCmd(ctx *impls.HandlerContext, msg *cproto.ReqSystemCmd, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	sid := ctx.GetSession().GetID()
	uid := convert.TryParseInt64(sid, 0)
	if uid == 0 {
		l.L.Errorf("cmd error uid")
		return
	}
	player := s.Mgr.GetPlayer(uid)
	player.SystemRequest(msg.System, msg.Cmd, msg.Args, func(ret []byte, code int32) {
		apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.AckSystemCmd{
			Ret:  ret,
			Code: code,
		})
	})
}
