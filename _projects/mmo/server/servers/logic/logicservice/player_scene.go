package logicservice

import (
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/client/session"
	common1 "github.com/dfklegend/cell2/utils/common"
	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/waterfall"

	"mmo/common/define"
	mymsg "mmo/messages"
	"mmo/messages/cproto"
)

// 场景相关行为

func (p *LogicPlayer) IsInScene() bool {
	return p.sceneId > 0
}

func (p *LogicPlayer) GetSceneServer() string {
	return p.sceneServer
}

func (p *LogicPlayer) GetSceneId() uint64 {
	return p.sceneId
}

// AllocEnterScene 创建一个关卡副本，并且进入场景
func (p *LogicPlayer) AllocEnterScene(cfgId int, pos define.Pos, notifyClient bool, switchLine bool, cb func(succ bool)) {
	// . 离开当前场景
	// . 请求分配战斗服
	// . 请求进入战斗服

	l.L.Errorf("AllocEnterScene: %v", cfgId)

	// 后续cfgId可以知道是什么战斗类型
	ns := p.ns
	var allocAck *mymsg.SMAllocSceneAck
	waterfall.NewBuilder(ns.GetRunService().GetScheduler()).
		Next(func(callback waterfall.Callback, args ...interface{}) {
			if p.IsInScene() {
				p.ReqLeaveCurScene(func(succ bool) {
					callback(!succ)
				})
				return
			}
			callback(false)
		}).
		Next(func(callback waterfall.Callback, args ...interface{}) {

			l.L.Infof("allocscene: %v", cfgId)
			app.Request(ns, "scenem.remote.allocscene", nil, &mymsg.SMAllocScene{
				UId:   p.uid,
				CfgId: int32(cfgId),
			}, func(err error, raw interface{}) {
				if err != nil {
					l.L.Errorf("scenem allocscene failed")
					callback(true)
					return
				}

				ack := raw.(*mymsg.SMAllocSceneAck)
				if ack == nil {
					l.L.Errorf("bad SMAllocSceneAck")
					callback(true)
				}

				allocAck = ack

				l.L.Infof("SMAllocSceneAck: %v", allocAck)
				//mode.InitSceneInfo(ack.ServiceId, ack.SceneId, ack.Token)
				callback(false)
			})
		}).
		Next(func(callback waterfall.Callback, args ...interface{}) {
			app.Request(ns, "scene.remote.enter", allocAck.ServiceId, &mymsg.SceneEnter{
				UId:        p.uid,
				SceneId:    allocAck.SceneId,
				Token:      allocAck.Token,
				FrontId:    p.FrontId,
				NetId:      p.NetId,
				LogicId:    ns.Name,
				CfgId:      int32(cfgId),
				SwitchLine: switchLine,
			}, func(err error, raw any) {
				if err != nil {
					l.L.Errorf("scene enter error: %v", err)
					callback(true)
					return
				}

				p.onEnterScene(notifyClient, allocAck.SceneId, allocAck.ServiceId, int32(cfgId))
				l.L.Infof("scene enter succ")
				callback(false)
			})

		}).
		Final(func(err bool, args ...interface{}) {
			cb(!err)
		}).
		Do()
}

func (p *LogicPlayer) ReEnterScene(cb func(succ bool)) {
	ns := p.ns
	app.Request(ns, "scene.remote.reenter", p.sceneServer, &mymsg.SceneEnter{
		UId:     p.uid,
		SceneId: p.sceneId,
		Token:   0, // not used
		FrontId: p.FrontId,
		NetId:   p.NetId,
		LogicId: ns.Name,
	}, func(err error, raw any) {
		if err != nil {
			l.L.Errorf("scene reenter error: %v", err)
			cb(false)
			return
		}

		l.L.Infof("scene reenter succ")
		p.bindScene(p.sceneServer)
		cb(true)
	})
}

func (p *LogicPlayer) ReqLeaveCurScene(cb func(succ bool)) {
	if !p.IsInScene() {
		cb(true)
		return
	}

	var sceneId = p.sceneId
	app.Request(p.ns, "scene.remote.leave", p.sceneServer, &mymsg.SceneLeave{
		UId:     p.uid,
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

func (p *LogicPlayer) OnLeaveScene(sceneId uint64) {
	if sceneId != p.sceneId {
		l.L.Errorf("OnLeaveScene sceneId mismatch, cur: %v  leave: %v", p.sceneId, sceneId)
		return
	}

	passed := common1.NowMs() - p.sceneEnterTime
	l.L.Infof("%v leave scene: %v, stay in scene: %v ms", p.uid, sceneId, passed)

	p.sceneCfgId = 0
	p.sceneId = 0
	p.sceneServer = ""

	// 更新session
	// 下线则不用
	if !p.IsLogouting() {
		// 解绑scene
		p.bindScene("")
	}
}

func (p *LogicPlayer) onEnterScene(notifyClient bool, sceneId uint64, sceneServer string, cfgId int32) {
	p.sceneCfgId = int(cfgId)
	p.sceneId = sceneId
	p.sceneServer = sceneServer
	p.sceneEnterTime = common1.NowMs()
	p.RefreshSceneKeepAlive()

	p.ns.GetLogger().Infof("%v enterscene: %v %v", p.uid, sceneServer, sceneId)

	p.bindScene(sceneServer)
	// 通知客户端进入新场景
	if notifyClient {
		p.notifyClientOnChangeScene(sceneServer, cfgId, sceneId)
	}
}

func (p *LogicPlayer) bindScene(sceneServer string) {
	// 绑定一下 scene服务
	bs := session.NewBackSession(p.ns, p.FrontId, p.NetId, "")
	bs.Set("scene", sceneServer)
	bs.PushSession(func(err error) {
		if err == nil {
			p.ns.GetLogger().Infof("bindScene %v succ", sceneServer)
		}
	})
}

func (p *LogicPlayer) notifyClientOnChangeScene(sceneServer string, cfgId int32, sceneId uint64) {
	p.ns.GetLogger().Infof("notifyClientOnChangeScene %v", sceneId)
	// 通知客户端场景改变了
	p.PushMsg("onchangescene", &cproto.LoadScene{
		ServerId: sceneServer,
		CfgId:    cfgId,
		SceneId:  sceneId,
	})
}

func (p *LogicPlayer) switchSceneLock() bool {
	if p.switchSceneLocked {
		return false
	}
	p.switchSceneLocked = true
	return true
}

func (p *LogicPlayer) switchSceneUnlock() {
	p.switchSceneLocked = false
}
