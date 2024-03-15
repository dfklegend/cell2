package sceneservice

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
	define2 "mmo/servers/scene/define"
	"mmo/servers/scene/sceneplayer/cmds/context"
	"mmo/utils"
)

func init() {
	registry.Registry.AddCollection("scene.handler").
		Register(&Handler{}, apientry.WithName("scene"), apientry.WithNameFunc(strings.ToLower))
}

type Handler struct {
	api.APIEntry
}

func (e *Handler) ReqCharInfo(ctx *impls.HandlerContext, msg *cproto.EmptyArg, cbFunc apientry.HandlerCBFunc) {
	// 请求角色信息，登录成功之后
	s := ctx.ActorContext.Actor().(*Service)
	logger := s.GetLogger()
	bs := ctx.GetBackSession()
	uid := convert.TryParseInt64(bs.GetID(), 0)

	logger.Infof("got req charinfo:%v ", uid)
	player := s.players.GetPlayer(uid)
	if uid == 0 || player == nil {
		logger.Infof("can not find player:%v ", uid)
		apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.NormalAck{
			Code: int32(define.ErrFaild),
		})
		return
	}

	// 请求数据
	player.ReqCharInfo()
	//
	apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.NormalAck{
		Code: int32(define.Succ),
	})
}

func (e *Handler) OpenCamera(ctx *impls.HandlerContext, msg *cproto.EmptyArg, cbFunc apientry.HandlerCBFunc) {
	// 打开摄像机(查看场景)
	s := ctx.ActorContext.Actor().(*Service)
	logger := s.GetLogger()
	bs := ctx.GetBackSession()
	uid := convert.TryParseInt64(bs.GetID(), 0)

	logger.Infof("got req charinfo:%v ", uid)
	player := s.players.GetPlayer(uid)
	if uid == 0 || player == nil {
		logger.Infof("can not find player:%v ", uid)
		apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.NormalAck{
			Code: int32(define.ErrFaild),
		})
		return
	}

	if player.GetState() != define2.Normal {
		logger.Infof("player:%v is incorrect state: %v", uid, player.GetState())
		apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.NormalAck{
			Code: int32(define.ErrFaild),
		})
		return
	}
	player.OpenCamera()

	apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.NormalAck{
		Code: int32(define.Succ),
	})
}

func (e *Handler) ClientLoadSceneOver(ctx *impls.HandlerContext, msg *cproto.ClientLoadSceneOver, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	logger := s.GetLogger()

	logger.Infof("ClientLoadSceneOver: %v %v", msg.UId, msg.SceneId)
	scene := s.Mgr.FindScene(msg.SceneId)
	if scene == nil {
		logger.Infof("ClientLoadSceneOver can not find scene: %v %v", msg.UId, msg.SceneId)
		apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.NormalAck{
			Code: int32(define.ErrFaild),
		})
		return
	}
	scene.OnClientSceneLoadOver(msg.UId)

	apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.NormalAck{
		Code: int32(define.Succ),
	})
}

// Cmd
// 提供了简单的字符串命令协议
func (e *Handler) Cmd(ctx *impls.HandlerContext, msg *cproto.Cmd, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	uid := utils.GetIdFromContext(ctx)
	player := s.GetPlayers().GetPlayer(uid)
	if uid == 0 || player == nil {
		l.L.Errorf("cmd error uid")
		apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.NormalAck{
			Code: int32(define.ErrFaild),
		})
		return
	}

	cmdCtx := context.NewCmdContext(player)
	s.cmdMgr.Dispatch(cmdCtx, msg.Cmd, func(result string) {
		apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.CmdAck{
			Result: result,
		})
	})
}

// SystemCmd
// 系统子协议
func (e *Handler) SystemCmd(ctx *impls.HandlerContext, msg *cproto.ReqSystemCmd, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	logger := s.GetLogger()
	uid := utils.GetIdFromContext(ctx)
	//logger.Infof("SystemCmd: %v", uid)
	player := s.GetPlayers().GetPlayer(uid)
	if uid == 0 || player == nil {
		logger.Errorf("systemCmd error uid")
		apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.NormalAck{
			Code: int32(define.ErrFaild),
		})
		return
	}

	player.SystemRequest(msg.System, msg.Cmd, msg.Args, func(ret []byte, code int32) {
		apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.AckSystemCmd{
			Ret:  ret,
			Code: code,
		})
	})
}
