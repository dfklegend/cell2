package logic

import (
	"strings"

	as "github.com/dfklegend/cell2/actorex/service"
	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/apimapper/registry"
	"github.com/dfklegend/cell2/node/app"
	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/waterfall"

	"mmo/common/define"
	mymsg "mmo/messages"
	"mmo/servers/logic_old"
)

func init() {
	registry.Registry.AddCollection("logic_old.remote").
		Register(&Entry{}, apientry.WithName("logicremote"), apientry.WithNameFunc(strings.ToLower))
}

func Visit() {
}

type Entry struct {
	api.APIEntry
}

// LoadPlayer 载入玩家
// LogicLoadPlayer
func (e *Entry) LoadPlayer(ctx *as.RemoteContext, msg *mymsg.LogicLoadPlayer, cbFunc apientry.HandlerCBFunc) {
	// . 检查是否已经有该角色
	// . 创建角色
	// . 请求db载入角色
	// . 通知center上线
	// . player.LoadInfo
	s := ctx.ActorContext.Actor().(*Service)

	l.L.Infof("%v LoadPlayer enter:%v", s.Name, msg.UId)
	if !s.State.CanWork() {
		l.L.Errorf("logic_old state can not work: %v", s.Name)
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
			Code: int32(define.ErrWithStr),
		})
		return
	}

	player := s.Mgr.GetPlayer(msg.UId)
	if player != nil {
		l.L.Errorf("already has a player: %v", msg.UId)
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
			Code: int32(define.ErrAlreadyOnline),
		})
		return
	}

	player = s.Mgr.CreatePlayer(msg.UId)
	player.Init(msg.ServerId, msg.NetId)

	ns := s.GetNodeService()

	//var info *mymsg.PlayerInfo
	//newCreatePlayer := false

	waterfall.NewBuilder(ns.GetRunService().GetScheduler()).
		Next(func(callback waterfall.Callback, args ...interface{}) {

			// ask handler to load player
			app.Request(ns, "db.dbremote.loadplayer", nil, &mymsg.DBLoadPlayer{
				UId: msg.UId,
			}, func(err error, raw interface{}) {
				if err != nil {
					l.L.Errorf("dbremote.loadplayer error: %v", err)
					callback(true)
					return
				}
				ack := raw.(*mymsg.DBLoadPlayerAck)
				if ack == nil {
					l.L.Errorf("dbremote.loadplayer ret is error ")
					callback(true)
					return
				}

				if ack.Info == nil && !ack.NewPlayer {
					l.L.Errorf("got nil info and not NewPlayer")
					callback(true)
					return
				}

				//newCreatePlayer = ack.NewPlayer
				//info = ack.Info
				callback(false)
			})

		}).
		Next(func(callback waterfall.Callback, args ...interface{}) {

			// notify handler
			//app.Request(ns, "center.centerremote.onlogiconline", nil, &mymsg.OnLogicOffline{
			//	UId:     msg.UId,
			//	LogicId: s.Name,
			//}, func(err error, raw interface{}) {
			//	if err != nil {
			//		l.L.Errorf("centerremote.onlogiconline error: %v", err)
			//		callback(true)
			//		return
			//	}
			//
			//	ack := raw.(*mymsg.NormalAck)
			//	if ack == nil || ack.Code != int32(define.Succ) {
			//		l.L.Errorf("centerremote.onlogiconline failed ")
			//		callback(true)
			//		return
			//	}
			//
			//	if newCreatePlayer {
			//		player.InitInfo()
			//	} else {
			//		player.LoadInfo(info)
			//	}
			//
			//	player.EnterWorld()
			//	callback(false)
			//})

		}).
		Final(func(err bool, args ...interface{}) {
			if err {
				apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
					Code: int32(define.ErrWithStr),
				})
			} else {
				apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
					Code: int32(define.Succ),
				})
			}
		}).
		Do()
}

// ReqOffline 登出 handler->
func (e *Entry) ReqOffline(ctx *as.RemoteContext, msg *mymsg.ReqOffline, cbFunc apientry.HandlerCBFunc) {

	s := ctx.ActorContext.Actor().(*Service)
	l.L.Infof("%v ReqOffline enter:%v", s.Name, msg.UId)

	//player := s.Mgr.GetPlayer(msg.UId)
	//if player == nil {
	//	l.L.Warnf("can not find player: %v", msg.UId)
	//	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{})
	//	return
	//}
	//
	//// 存在风险: 下线过程中，如果有获取物品，将丢失
	//
	//ns := s.GetNodeService()
	//waterfall.NewBuilder(ns.GetRunService().GetScheduler()).
	//	Next(func(callback waterfall.Callback, args ...interface{}) {
	//
	//		//callback := args[0].(waterfall.Callback)
	//		if player.IsInScene() {
	//			player.NotifySceneOffline()
	//			callback(false)
	//		} else {
	//			callback(false)
	//		}
	//
	//	}).
	//	Next(func(callback waterfall.Callback, args ...interface{}) {
	//
	//		// saveplayer
	//		if !player.IsDirt() {
	//			callback(false)
	//			return
	//		}
	//
	//		player.ClearDirt()
	//		app.Request(ns, "handler.dbremote.saveplayer", nil, &mymsg.DBSavePlayer{
	//			Info: player.MakeData(),
	//		}, func(err error, raw interface{}) {
	//			if err != nil {
	//				l.L.Errorf("dbremote.saveplayer error: %v", err)
	//				callback(true)
	//				return
	//			}
	//			ack := raw.(*mymsg.NormalAck)
	//			if ack == nil || ack.Code != int32(define.Succ) {
	//				l.L.Errorf("dbremote.saveplayer error ")
	//				callback(true)
	//				return
	//			}
	//
	//			// 后面考虑一种状态，会再次存储，存储之后自动offline
	//			if player.IsDirt() {
	//				l.L.Warnf("player %v dirt again!", player.UId)
	//			}
	//			callback(false)
	//		})
	//
	//	}).
	//	Next(func(callback waterfall.Callback, args ...interface{}) {
	//
	//		//
	//		app.Request(ns, "handler.centerremote.onlogicoffline", nil,
	//			&mymsg.OnLogicOffline{
	//				UId: msg.UId,
	//			}, func(err error, raw interface{}) {
	//				callback(false)
	//			})
	//
	//	}).
	//	Next(func(callback waterfall.Callback, args ...interface{}) {
	//
	//		//
	//		s.Mgr.DestroyPlayer(msg.UId)
	//		callback(false)
	//
	//	}).
	//	Final(func(err bool, args ...interface{}) {
	//		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{})
	//	}).
	//	Do()

	logic_old.ReqOffline(s.GetNodeService(), s.Mgr, msg.UId, func(err bool) {
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{})
	})
}

func (e *Entry) OnLeaveScene(ctx *as.RemoteContext, msg *mymsg.OnLeaveScene, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	l.L.Infof("%v OnLeaveScene enter:%v", s.Name, msg.UId)

	player := s.Mgr.GetPlayer(msg.UId)
	if player == nil {
		l.L.Errorf("OnLeaveScene can not find player: %v", msg.UId)
		return
	}

	player.OnLeaveScene(msg.SceneId)
	apientry.CheckInvokeCBFunc(cbFunc, nil, nil)
}

func (e *Entry) OnFightResult(ctx *as.RemoteContext, msg *mymsg.OnFightResult, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	l.L.Infof("%v OnFightResult enter:%v", s.Name, msg.UId)

	player := s.Mgr.GetPlayer(msg.UId)
	if player == nil {
		l.L.Errorf("OnFightResult can not find player: %v", msg.UId)
		return
	}

	player.OnGotFightReward(int(msg.Money), int(msg.Exp))
	apientry.CheckInvokeCBFunc(cbFunc, nil, nil)
}

func (e *Entry) OnCardFightResult(ctx *as.RemoteContext, msg *mymsg.SceneCardResult, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	l.L.Infof("%v OnFightResult enter:%v", s.Name, msg.UId)

	player := s.Mgr.GetPlayer(msg.UId)
	if player == nil {
		l.L.Errorf("OnFightResult can not find player: %v", msg.UId)
		return
	}

	player.GetCurFight().OnFightResult(msg)
	apientry.CheckInvokeCBFunc(cbFunc, nil, nil)
}
