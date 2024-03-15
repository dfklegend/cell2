package sceneservice

import (
	"errors"
	"strings"

	as "github.com/dfklegend/cell2/actorex/service"
	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/apimapper/registry"
	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/waterfall"

	"mmo/common/consts"
	"mmo/common/define"
	mymsg "mmo/messages"
	define2 "mmo/servers/scene/define"
)

func init() {
	registry.Registry.AddCollection("scene.remote").
		Register(&Entry{}, apientry.WithName("remote"), apientry.WithNameFunc(strings.ToLower))
}

func Visit() {
}

type Entry struct {
	api.APIEntry
}

func (e *Entry) AllocScene(ctx *as.RemoteContext, msg *mymsg.SAllocScene, cbFunc apientry.HandlerCBFunc) {
	// . 创建一个scene

	s := ctx.ActorContext.Actor().(*Service)
	l.L.Infof("%v AllocScene enter:%v %v %v", s.Name, msg.CfgId, msg.SceneId, msg.Token)

	succ := s.Mgr.AllocScene(msg.CfgId, msg.SceneId, int(msg.Token))
	if succ {
		apientry.CheckInvokeCBFunc(cbFunc, nil, nil)
	} else {
		apientry.CheckInvokeCBFunc(cbFunc, errors.New("alloc scene failed"), nil)
	}
}

// Enter
// 玩家进入场景服(上线或者切线)
// ack nil
func (e *Entry) Enter(ctx *as.RemoteContext, msg *mymsg.SceneEnter, cbFunc apientry.HandlerCBFunc) {
	e.doEnter(ctx, msg, cbFunc)
}

func (e *Entry) doEnter(ctx *as.RemoteContext, msg *mymsg.SceneEnter, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	logger := s.GetLogger()
	scene := s.Mgr.FindScene(msg.SceneId)
	switchLineEnter := msg.SwitchLine

	if scene == nil {
		logger.Errorf("can not find scene: %v", msg.SceneId)
		apientry.CheckInvokeCBFunc(cbFunc, consts.ErrCanNotFindScene, nil)
		return
	}

	if !scene.CanEnter(int(msg.Token)) {
		logger.Errorf("error token: %v", msg.Token)
		apientry.CheckInvokeCBFunc(cbFunc, consts.ErrNotImportant, nil)
		return
	}

	player := s.players.GetPlayer(msg.UId)
	if player != nil {
		logger.Errorf("scene %v alread has player %v", msg.SceneId, msg.UId)
		apientry.CheckInvokeCBFunc(cbFunc, consts.ErrNotImportant, nil)
		return
	}

	player = s.players.CreatePlayer(msg.UId, msg.FrontId, msg.NetId, msg.LogicId)
	ns := s.NodeService
	waterfall.NewBuilder(ns.GetRunService().GetScheduler()).
		Next(func(callback waterfall.Callback, args ...interface{}) {
			// 玩家载入数据，初始化
			player.Enter(switchLineEnter, func(succ bool) {
				callback(!succ)
			})
		}).
		Next(func(callback waterfall.Callback, args ...interface{}) {

			if switchLineEnter && define.DebugForceSwitchLineSceneEnterFail1 {
				callback(true)
				return
			}

			// 进入场景
			if !scene.PlayerEnter(player) {
				callback(true)
			} else {
				if switchLineEnter && define.DebugForceSwitchLineSceneEnterFail2 {
					callback(true)
					return
				}
				callback(false)
			}
		}).
		Final(func(err bool, args ...any) {
			if err {
				// player设置成等待删除，不要触发存储
				player.ChangeState(define2.WaitMgrDelete)
				apientry.CheckInvokeCBFunc(cbFunc, consts.ErrToDefine, nil)
			} else {
				apientry.CheckInvokeCBFunc(cbFunc, nil, nil)
				player.EnterPreNormal()
			}
		}).
		Do()
}

func (e *Entry) ReEnter(ctx *as.RemoteContext, msg *mymsg.SceneEnter, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	player := s.players.GetPlayer(msg.UId)

	if player == nil {
		l.L.Errorf("scene %v has not player %v", msg.SceneId, msg.UId)
		apientry.CheckInvokeCBFunc(cbFunc, consts.ErrToDefine, nil)
		return
	}

	scene := s.Mgr.FindScene(msg.SceneId)
	if scene == nil {
		l.L.Errorf("can not find scene %v ", msg.SceneId)
		apientry.CheckInvokeCBFunc(cbFunc, consts.ErrToDefine, nil)
		return
	}

	player.ReEnter(msg.FrontId, msg.NetId)
	//player.ChangeScene(scene, msg.CfgId, msg.SceneId)
	apientry.CheckInvokeCBFunc(cbFunc, nil, nil)
}

// PreCheckPermit 提前检查是否允许进入
// ack error(nil), nil
func (e *Entry) PreCheckPermit(ctx *as.RemoteContext, msg *mymsg.SceneEnter, cbFunc apientry.HandlerCBFunc) {
	// 检查合法性
	s := ctx.ActorContext.Actor().(*Service)
	logger := s.GetLogger()
	scene := s.Mgr.FindScene(msg.SceneId)

	if scene == nil {
		logger.Errorf("can not find scene: %v", msg.SceneId)
		apientry.CheckInvokeCBFunc(cbFunc, consts.ErrCanNotFindScene, nil)
		return
	}

	if !scene.CanEnter(int(msg.Token)) {
		logger.Errorf("error token: %v", msg.Token)
		apientry.CheckInvokeCBFunc(cbFunc, consts.ErrNotImportant, nil)
		return
	}

	// 如果需要，后续可以加个防止场景后续很短时间被释放的功能
	apientry.CheckInvokeCBFunc(cbFunc, nil, nil)
}

//func (e *Entry) LogicOffline(ctx *as.RemoteContext, msg *mymsg.LogicOffline) {
//	s := ctx.ActorContext.Actor().(*Service)
//	scene := s.Mgr.FindScene(msg.SceneId)
//	if scene == nil {
//		l.L.Errorf("can not find scene: %v", msg.SceneId)
//		return
//	}
//
//	scene.PlayerLeave(msg.UId)
//}

// Leave
// 切线时，会被要求先离开
func (e *Entry) Leave(ctx *as.RemoteContext, msg *mymsg.SceneLeave, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	logger := s.GetLogger()
	scene := s.Mgr.FindScene(msg.SceneId)
	if scene == nil {
		logger.Errorf("can not find scene: %v", msg.SceneId)
		apientry.CheckInvokeCBFunc(cbFunc, consts.ErrCanNotFindScene, nil)
		return
	}

	player := s.players.GetPlayer(msg.UId)
	if player == nil {
		logger.Errorf("can not find player: %v", msg.UId)
		apientry.CheckInvokeCBFunc(cbFunc, consts.ErrToDefine, nil)
		return
	}

	if player.GetSceneId() != msg.SceneId {
		logger.Errorf("player: %v is not in scene: %v, curScene: %v", msg.UId, msg.SceneId, player.GetSceneId())
		apientry.CheckInvokeCBFunc(cbFunc, consts.ErrToDefine, nil)
		return
	}

	player.Leave(func(succ bool) {
		code := define.Succ
		if !succ {
			code = define.ErrFaild
		}
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
			Code: int32(code),
		})
	})
}

func (e *Entry) ChangeScene(ctx *as.RemoteContext, msg *mymsg.ChangeScene, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	logger := s.GetLogger()

	logger.Infof("changescene enter: %v %v", msg.UId, msg.SceneId)
	scene := s.Mgr.FindScene(msg.SceneId)
	if scene == nil {
		logger.Errorf("can not find scene: %v", msg.SceneId)
		apientry.CheckInvokeCBFunc(cbFunc, consts.ErrCanNotFindScene, nil)
		return
	}

	player := s.players.GetPlayer(msg.UId)
	if player == nil {
		logger.Errorf("can not find player: %v", msg.UId)
		apientry.CheckInvokeCBFunc(cbFunc, consts.ErrToDefine, nil)
		return
	}

	if !scene.CheckToken(int(msg.Token)) {
		logger.Errorf("token mismatch", msg.Token)
		apientry.CheckInvokeCBFunc(cbFunc, consts.ErrToDefine, nil)
		return
	}

	player.LeaveCurScene()
	code := define.Succ
	if !scene.PlayerEnter(player) {
		code = define.ErrFaild
	}

	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
		Code: int32(code),
	})
}

func (e *Entry) OnOffline(ctx *as.RemoteContext, msg *mymsg.OnOffline, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	player := s.players.GetPlayer(msg.UId)
	logger := s.GetLogger()
	if player == nil {
		logger.Errorf("OnOffline can not find player: %v", msg.UId)
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
			Code: int32(define.ErrFaild),
		})
		return
	}

	player.OnOffline()
	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
		Code: int32(define.Succ),
	})
}

func (e *Entry) OnReOnline(ctx *as.RemoteContext, msg *mymsg.OnReOnline, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	player := s.players.GetPlayer(msg.UId)
	if player == nil {
		l.L.Errorf("OnReOnline can not find player: %v", msg.UId)
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
			Code: int32(define.ErrFaild),
		})
		return
	}

	player.OnReOnline()
	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
		Code: int32(define.Succ),
	})
}

func (e *Entry) InitCardFight(ctx *as.RemoteContext, msg *mymsg.SceneInitCardFight, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	scene := s.Mgr.FindScene(msg.SceneId)
	if scene == nil || !scene.CheckToken(int(msg.Token)) {
		l.L.Errorf("can not find scene or token mismatch: %v", msg.SceneId)
		apientry.CheckInvokeCBFunc(cbFunc, consts.ErrCanNotFindScene, nil)
		return
	}

	scene.SetInitData(msg)
	apientry.CheckInvokeCBFunc(cbFunc, nil, nil)
}
