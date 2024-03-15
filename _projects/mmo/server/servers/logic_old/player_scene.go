package logic_old

import (
	"fmt"

	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/client/session"
	"github.com/dfklegend/cell2/utils/common"
	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/waterfall"

	mymsg "mmo/messages"
	"mmo/messages/cproto"
)

// LeaveScene
// AllocEnterScene
// SwitchScene

func (p *Player) IsInScene() bool {
	return p.sceneId > 0
}

func (p *Player) QueryAndEnterScene(sceneId uint64, cb func(succ bool)) {
	// . 取到场景token
	// . 请求进入
	waterfall.NewBuilder(p.ns.GetRunService().GetScheduler()).
		Next(func(callback waterfall.Callback, args ...any) {
			app.Request(p.ns, "scenem.remote.queryscene", nil, &mymsg.SMQueryScene{
				SceneId: sceneId,
			}, func(err error, ret any) {
				if err != nil {
					l.L.Errorf("query scene failed, err: %v", err)
					callback(true)
					return
				}
				callback(false, ret)
			})
		}).
		Next(func(callback waterfall.Callback, args ...any) {
			ack := args[0].(*mymsg.SMQuerySceneAck)
			p.ReqEnterScene(ack.ServiceId, ack.SceneId, ack.Token, ack.CfgId, func(succ bool) {
				callback(!succ)
			})
		}).
		Final(func(err bool, args ...any) {
			if err {
				l.L.Errorf("QueryAndEnterScene failed")
			}
			cb(!err)
		}).Do()
}

func (p *Player) ReqEnterScene(sceneServiceId string, sceneId uint64, token int32, cfgId int32, cb func(succ bool)) {

	waterfall.NewBuilder(p.ns.GetRunService().GetScheduler()).
		Next(func(callback waterfall.Callback, args ...any) {
			if p.IsInScene() {
				// 首先得离开老场景
				p.ReqLeaveCurScene(func(succ bool) {
					callback(!succ)
				})
			} else {
				callback(false)
			}
		}).
		Next(func(callback waterfall.Callback, args ...any) {
			app.Request(p.ns, "scene.remote.enter", sceneServiceId, &mymsg.SceneEnter{
				UId:     p.UId,
				SceneId: sceneId,
				Token:   token,
				FrontId: p.FrontId,
				NetId:   p.NetId,
				LogicId: p.ns.Name,
			}, func(err error, raw any) {
				if err != nil {
					l.L.Errorf("scene enter error: %v", err)
					cb(false)
					return
				}

				p.onEnterScene(sceneId, sceneServiceId, cfgId)
				l.L.Infof("scene enter succ")
				cb(true)
			})
		}).
		Final(func(err bool, args ...any) {

		}).Do()

}

func (p *Player) onEnterScene(sceneId uint64, sceneServer string, cfgId int32) {
	p.sceneId = sceneId
	p.sceneServer = sceneServer
	p.sceneEnterTime = common.NowMs()

	l.L.Infof("%v enterscene: %v %v", p.UId, sceneServer, sceneId)

	p.PushBattleLog(fmt.Sprintf("进入战斗，server: %v scene: %v", sceneServer, sceneId))

	// 绑定一下 scene服务
	bs := session.NewBackSession(p.ns, p.FrontId, p.NetId, "")
	bs.Set("scene", sceneServer)
	bs.PushSession(func(err error) {
		if err == nil {
			//l.L.Infof("push session succ")
			p.notifyClientLoadScene(int(cfgId), sceneId)
		}
	})
}

// notifyClientLoadScene 要求客户端开始载入场景
func (p *Player) notifyClientLoadScene(cfgId int, sceneId uint64) {
	p.PushMsg("loadscene", &cproto.LoadScene{
		CfgId:   int32(cfgId),
		SceneId: sceneId,
	})
}

func (p *Player) ReqLeaveCurScene(cb func(succ bool)) {
	if !p.IsInScene() {
		cb(true)
		return
	}

	var sceneId = p.sceneId
	app.Request(p.ns, "scene.remote.leave", p.sceneServer, &mymsg.SceneLeave{
		UId:     p.UId,
		SceneId: sceneId,
	}, func(err error, ret any) {
		if err != nil {
			l.L.Errorf("scene.remote.leave error: %v", err)
		} else {
			p.OnLeaveScene(sceneId)
		}
		cb(err == nil)
	})
}

func (p *Player) OnLeaveScene(sceneId uint64) {
	if sceneId != p.sceneId {
		l.L.Errorf("OnLeaveScene sceneId mismatch, cur: %v  leave: %v", p.sceneId, sceneId)
		return
	}

	passed := common.NowMs() - p.sceneEnterTime
	l.L.Infof("%v leave scene: %v, stay in scene: %v ms", p.UId, sceneId, passed)

	p.sceneId = 0
	p.sceneServer = ""

	// 检查是否需要, 有可能是下线造成的
	if !p.IsLogouting() {
		// 解绑scene
		bs := session.NewBackSession(p.ns, p.FrontId, p.NetId, "")
		bs.Set("scene", "")
		bs.PushSession(nil)
	}
}

// NotifySceneOffline 通知场景角色下线
func (p *Player) NotifySceneOffline() {
	p.ReqLeaveCurScene(func(succ bool) {})
}
