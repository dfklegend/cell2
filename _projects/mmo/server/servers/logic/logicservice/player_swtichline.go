package logicservice

import (
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/utils/waterfall"

	"mmo/common/define"
	mymsg "mmo/messages"
	"mmo/utils"
)

// TODO: 需要有个标志，保证logic上切线和切场景的独占性

// SwitchLine 切线
func (p *LogicPlayer) SwitchLine(sceneServer string, cfgId int32, sceneId uint64,
	token int32, pos define.Pos,
	cb func(succ bool)) {
	//. 先问问目标场景是否合法
	//. 向center请求切线(获得许可)
	//. 要求老的scene下线
	//. 要求在新的scene上线

	ns := p.ns
	logger := ns.GetLogger()

	if p.sceneId == sceneId {
		cb(true)
		return
	}

	if p.sceneServer == sceneServer {
		// 同一个服务，直接切换场景就行
		p.changeScene(cfgId, sceneId, token, cb)
		return
	}

	if define.DebugForceSwitchLineSceneEnterTokenErr {
		token = 0
	}

	needNotifySwitchLine := false

	logger.Infof("logic player switchline enter: %v  (%v, %v) -> (%v, %v)", p.uid,
		p.sceneServer, p.sceneId, sceneServer, sceneId)

	if !p.switchSceneLock() {
		logger.Infof("logic player SwitchLine failed, switchSceneLock failed, %v", p.uid)
		cb(false)
		return
	}

	waterfall.NewBuilder(ns.GetRunService().GetScheduler()).
		Next(func(callback waterfall.Callback, args ...interface{}) {
			app.Request(ns, "scene.remote.precheckpermit", sceneServer,
				&mymsg.SceneEnter{
					UId:     p.uid,
					SceneId: sceneId,
					Token:   token,
					FrontId: p.FrontId,
					NetId:   p.NetId,
					LogicId: ns.Name,
					CfgId:   cfgId,
				},
				func(err error, ret any) {
					if err != nil {
						logger.Errorf("scene precheckpermit failed")
						callback(true)
					} else {
						callback(false)
					}
				})
		}).
		Next(func(callback waterfall.Callback, args ...interface{}) {
			app.Request(ns, "center.centerremote.reqswitchline", nil,
				&mymsg.ReqSwitchLine{
					UId: p.uid,
				},
				func(err error, ret any) {
					ack := ret.(*mymsg.NormalAck)
					if err != nil || ack.Code != int32(define.Succ) {
						logger.Errorf("reqswitchline failed")
						callback(true)
					} else {
						callback(false)
						needNotifySwitchLine = true
					}
				})
		}).
		// 是否需要添加一个提前判定，目标场景存在与否，并且向目标场景预先申请下
		Next(func(callback waterfall.Callback, args ...interface{}) {
			app.Request(ns, "scene.remote.leave", p.sceneServer,
				&mymsg.SceneLeave{
					UId:     p.uid,
					SceneId: p.sceneId,
				},
				func(err error, ret any) {
					ack := ret.(*mymsg.NormalAck)
					if err != nil || ack.Code != int32(define.Succ) {
						logger.Errorf("scene leave failed: %v", p.uid)
						callback(true)
					} else {
						p.OnLeaveScene(p.sceneId)
						callback(false)
					}
				})
		}).
		Next(func(callback waterfall.Callback, args ...interface{}) {
			app.Request(ns, "scene.remote.enter", sceneServer,
				&mymsg.SceneEnter{
					UId:        p.uid,
					SceneId:    sceneId,
					Token:      token,
					FrontId:    p.FrontId,
					NetId:      p.NetId,
					LogicId:    ns.Name,
					CfgId:      cfgId,
					SwitchLine: true,
				},
				func(err error, ret any) {
					if err != nil {
						logger.Errorf("scene enter failed")
						callback(true)
					} else {

						p.onEnterScene(true, sceneId, sceneServer, int32(cfgId))
						callback(false)
					}
				})
		}).
		Final(func(err bool, args ...interface{}) {
			p.switchSceneUnlock()
			p.onSwitchLineEnd(!err, needNotifySwitchLine)
			cb(!err)
		}).
		Do()
}

// 通知
func (p *LogicPlayer) notifySwitchLineEnd(succ bool) {
	ns := p.ns
	app.Request(ns, "center.centerremote.onswitchlineend", nil,
		&mymsg.OnSwitchLineEnd{
			UId:  p.uid,
			Succ: succ,
		},
		func(err error, ret any) {
		})
}

func (p *LogicPlayer) onSwitchLineEnd(succ bool, needNotifySwitchLineEnd bool) {
	if needNotifySwitchLineEnd {
		p.notifySwitchLineEnd(succ)
	}
	if !succ {
		p.onSwitchLineFailed()
	}
}

func (p *LogicPlayer) onSwitchLineFailed() {
	p.goToSafeSceneIfNotInScene()
}

func (p *LogicPlayer) changeScene(cfgId int32, sceneId uint64, token int32, cb func(succ bool)) {
	ns := p.ns
	logger := ns.GetLogger()
	sceneServer := p.sceneServer

	logger.Infof("logic player changescene enter: %v (%v) -> (%v)", p.uid, p.sceneId, sceneId)
	if !p.switchSceneLock() {
		logger.Infof("logic player changescene failed, switchSceneLock failed, %v", p.uid)
		cb(false)
		return
	}

	waterfall.NewBuilder(ns.GetRunService().GetScheduler()).
		Next(func(callback waterfall.Callback, args ...interface{}) {
			app.Request(ns, "scene.remote.precheckpermit", sceneServer,
				&mymsg.SceneEnter{
					UId:     p.uid,
					SceneId: sceneId,
					Token:   token,
					FrontId: p.FrontId,
					NetId:   p.NetId,
					LogicId: ns.Name,
					CfgId:   cfgId,
				},
				func(err error, ret any) {
					if err != nil {
						logger.Errorf("scene precheckpermit failed")
						callback(true)
					} else {
						callback(false)
					}
				})
		}).
		Next(func(callback waterfall.Callback, args ...interface{}) {
			app.Request(ns, "scene.remote.changescene", sceneServer,
				&mymsg.ChangeScene{
					UId:     p.uid,
					SceneId: sceneId,
					Token:   token,
					CfgId:   cfgId,
				},
				func(err error, ret any) {
					ack := utils.TryGetNormalAck(ret)
					if err != nil || (ack != nil && ack.Code != int32(define.Succ)) {
						logger.Errorf("changescene enter failed")
						p.onChangeSceneFailed()
						callback(true)
					} else {
						p.onEnterScene(true, sceneId, sceneServer, int32(cfgId))
						callback(false)
					}
				})
		}).
		Final(func(err bool, args ...interface{}) {
			p.switchSceneUnlock()
			cb(!err)
		}).
		Do()
}

func (p *LogicPlayer) onChangeSceneFailed() {
	p.goToSafeSceneIfNotInScene()
}

func (p *LogicPlayer) goToSafeSceneIfNotInScene() {
	// 检查一下，是否在场景中
	if p.IsInScene() {
		return
	}
	// 将玩家送到某个场景里
	p.goToSafeScene(nil)
}

// 各种场景异常之后，都可以进入一个绝对不会失败的家园场景
func (p *LogicPlayer) goToSafeScene(cb func(succ bool)) {
	// 思考: 如果是mmo,可能应该是存储的最后安全地点，比如副本，那就是副本入口
	// 基本不可能失败
	p.AllocEnterScene(int(define.ScenePlayerHome), define.Pos{}, true, true, func(succ bool) {
		if cb != nil {
			cb(succ)
		}
	})
}
