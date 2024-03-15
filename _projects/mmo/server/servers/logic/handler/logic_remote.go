package handler

import (
	"strings"

	as "github.com/dfklegend/cell2/actorex/service"
	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/apimapper/registry"
	"github.com/dfklegend/cell2/node/app"
	l "github.com/dfklegend/cell2/utils/logger"

	"mmo/common/define"
	mymsg "mmo/messages"
	logicservice "mmo/servers/logic/logicservice"
)

func init() {
	registry.Registry.AddCollection("logic.remote").
		Register(&Entry{}, apientry.WithName("logicremote"), apientry.WithNameFunc(strings.ToLower))
}

func Visit() {
}

type Entry struct {
	api.APIEntry
}

// Enter
// 进入游戏
// callby gate
func (e *Entry) Enter(ctx *as.RemoteContext, msg *mymsg.LogicLoadPlayer, cbFunc apientry.HandlerCBFunc) {
	// . 检查是否已经有该角色
	// . 创建角色
	// . 请求db载入角色
	// . 通知center上线
	// . player.LoadInfo
	// . 创建场景对象
	s := ctx.ActorContext.Actor().(*logicservice.Service)
	ns := s.NodeService

	l.L.Infof("%v logicremote.enter enter:%v", s.Name, msg.UId)
	players := s.GetPlayers()
	player := players.GetPlayer(msg.UId)
	if player != nil {
		ns.GetLogger().Errorf("already has a player: %v", msg.UId)
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
			Code: int32(define.ErrAlreadyOnline),
		})
		return
	}

	player = players.CreatePlayer(msg.UId)
	player.Init(msg.ServerId, msg.NetId)
	player.Enter(msg, cbFunc)
}

func (e *Entry) ReEnter(ctx *as.RemoteContext, msg *mymsg.LogicLoadPlayer, cbFunc apientry.HandlerCBFunc) {
	// .
	s := ctx.ActorContext.Actor().(*logicservice.Service)

	l.L.Infof("%v logicremote.reenter enter:%v", s.Name, msg.UId)
	players := s.GetPlayers()
	player := players.GetPlayer(msg.UId)
	if player == nil {
		l.L.Errorf("can not find player: %v", msg.UId)
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
			Code: int32(define.ErrFaild),
		})
		return
	}

	player.Init(msg.ServerId, msg.NetId)
	player.ReEnter(cbFunc)
}

// OnOffline 客户端连接断开
func (e *Entry) OnOffline(ctx *as.RemoteContext, msg *mymsg.OnOffline, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*logicservice.Service)

	l.L.Infof("%v logicremote.OnOffline enter:%v", s.Name, msg.UId)
	players := s.GetPlayers()
	player := players.GetPlayer(msg.UId)
	if player == nil {
		l.L.Errorf("can not find player: %v", msg.UId)
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
			Code: int32(define.ErrFaild),
		})
		return
	}

	if !player.IsInScene() {
		// 不可能出现
		l.L.Errorf("OnOffline critical error player: %v not in scene", msg.UId)
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
			Code: int32(define.ErrFaild),
		})
		return
	}

	player.OnOffline()
	// 通知scene
	ns := s.GetNodeService()
	app.Request(ns, "scene.remote.onoffline", player.GetSceneServer(), msg, func(err error, ack any) {
		if err != nil {
			apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
				Code: int32(define.ErrFaild),
			})
			return
		}
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
			Code: int32(define.Succ),
		})
	})
}

// OnReOnline
// 断线重连
func (e *Entry) OnReOnline(ctx *as.RemoteContext, msg *mymsg.OnReOnline, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*logicservice.Service)

	l.L.Infof("%v logicremote.OnReOnline enter:%v", s.Name, msg.UId)
	players := s.GetPlayers()
	player := players.GetPlayer(msg.UId)
	if player == nil {
		l.L.Errorf("can not find player: %v", msg.UId)
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
			Code: int32(define.ErrFaild),
		})
		return
	}

	if !player.IsInScene() {
		// 不可能出现
		l.L.Errorf("OnReOnline critical error player: %v not in scene", msg.UId)
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
			Code: int32(define.ErrFaild),
		})
		return
	}

	player.OnReOnline()

	// 通知scene
	ns := s.GetNodeService()
	app.Request(ns, "scene.remote.onreonline", player.GetSceneServer(), msg, func(err error, ack any) {
		if err != nil {
			apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
				Code: int32(define.ErrFaild),
			})
			return
		}
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
			Code: int32(define.Succ),
		})
	})
}

// ReqLogout
// scene ->
// 角色掉线时间久了，登出
func (e *Entry) ReqLogout(ctx *as.RemoteContext, msg *mymsg.ReqLogout, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*logicservice.Service)
	ns := s.GetNodeService()
	logger := ns.GetLogger()

	players := s.GetPlayers()
	player := players.GetPlayer(msg.UId)
	if player == nil {
		logger.Errorf("can not find player: %v", msg.UId)
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
			Code: int32(define.ErrFaild),
		})
		return
	}

	// 检查player是否有理由拒绝logout

	app.Request(ns, "center.centerremote.reqlogout", nil, msg, func(err error, ack any) {
		if err != nil {
			apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
				Code: int32(define.ErrFaild),
			})
			return
		}
		ack1 := ack.(*mymsg.NormalAck)
		if ack1.Code != int32(define.Succ) {
			// 处理一下
			logger.Infof("centerremote.reqlogout failed : %v", ack1.Code)
		}
		// 直接转发
		apientry.CheckInvokeCBFunc(cbFunc, nil, ack)
	})
}

// OnLogout
// scene ->
// scene通知登出完毕
func (e *Entry) OnLogout(ctx *as.RemoteContext, msg *mymsg.OnLogout, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*logicservice.Service)
	ns := s.GetNodeService()
	logger := ns.GetLogger()

	players := s.GetPlayers()
	player := players.GetPlayer(msg.UId)
	if player == nil {
		logger.Errorf("can not find player: %v", msg.UId)
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
			Code: int32(define.ErrFaild),
		})
		return
	}

	player.OnLogout()
	s.GetPlayers().DestroyPlayer(msg.UId)

	app.Request(ns, "center.centerremote.onlogout", nil, msg, func(err error, ack any) {
		apientry.CheckInvokeCBFunc(cbFunc, nil, ack)
	})
}

// ReqSwitchLine
// 处理切线请求
func (e *Entry) ReqSwitchLine(ctx *as.RemoteContext, msg *mymsg.LogicReqSwitchLine, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*logicservice.Service)
	ns := s.GetNodeService()
	logger := ns.GetLogger()

	players := s.GetPlayers()
	player := players.GetPlayer(msg.UId)
	if player == nil {
		logger.Errorf("can not find player: %v", msg.UId)
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
			Code: int32(define.ErrFaild),
		})
		return
	}

	player.SwitchLine(msg.SceneServer, msg.CfgId, msg.SceneId, msg.Token,
		define.Pos{X: msg.PosX, Y: msg.PosY, Z: msg.PosZ},
		func(succ bool) {
			code := define.Succ
			if !succ {
				code = define.ErrFaild
			}
			apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
				Code: int32(code),
			})
		})
}

func (e *Entry) CheckPlayer2(ctx *as.RemoteContext, msg *mymsg.CheckPlayer2, cbFunc apientry.HandlerCBFunc) {
	e.doKeepAlive(ctx, msg, cbFunc)
}

func (e *Entry) KeepAlive(ctx *as.RemoteContext, msg *mymsg.CheckPlayer2, cbFunc apientry.HandlerCBFunc) {
	e.doKeepAlive(ctx, msg, cbFunc)
}

func (e *Entry) doKeepAlive(ctx *as.RemoteContext, msg *mymsg.CheckPlayer2, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*logicservice.Service)
	ns := s.GetNodeService()
	logger := ns.GetLogger()

	players := s.GetPlayers()
	player := players.GetPlayer(msg.UId)
	if player == nil {
		logger.Errorf("doKeepAlive can not find player: %v", msg.UId)
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
			Code: int32(define.ErrFaild),
		})
		return
	}

	if player.GetSceneId() != msg.SceneId {
		logger.Errorf("doKeepAlive player %v scene mismatch: %v != %v", msg.UId, player.GetSceneId(), msg.SceneId)
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
			Code: int32(define.ErrFaild),
		})
		return
	}

	player.RefreshSceneKeepAlive()
	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
		Code: int32(define.Succ),
	})
}
