package sceneplayer

import (
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/utils/common"
	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/waterfall"

	mymsg "mmo/messages"
	"mmo/servers/scene/define"
)

// CanChangeScene
// 看看当前是否允许切换场景
func (p *ScenePlayer) CanChangeScene(tarSceneServer string) bool {
	if p.sceneServer != tarSceneServer {
		// 切线
		// 判断一下，比如如果支付过程中，不要切线
		return p.CanSwitchLine()
	}
	return true
}

func (p *ScenePlayer) CanSwitchLine() bool {
	return true
}

// ExitPointReqChangeScene
// 处理切换场景的请求
func (p *ScenePlayer) ExitPointReqChangeScene(tarCfgId int, pos define.Pos) {
	// . 判断是否允许
	// . 向scenem请求一个线(cfgId)
	// . 向logic请求切线
	logger := p.ns.GetLogger()

	if p.sceneSwitching {
		return
	}

	// 一些其他的合法性判定

	now := common.NowMs()
	if now < p.lastExitTriggled+1500 {
		return
	}

	p.lastExitTriggled = now
	p.sceneSwitching = true

	//
	ns := p.ns
	waterfall.NewBuilder(ns.GetRunService().GetScheduler()).
		Next(func(callback waterfall.Callback, args ...any) {
			app.Request(ns, "scenem.remote.reqscenebycfgid", nil,
				&mymsg.SMReqSceneByCfgId{
					CfgId: int32(tarCfgId),
				},
				func(err error, ret any) {
					if err != nil {
						logger.Infof("reqscenebycfgId failed: %v", err)
						callback(true)
						return
					}
					_, ok := ret.(*mymsg.SMQuerySceneAck)
					if !ok {
						logger.Infof("reqscenebycfgId failed, ack mismatch")
						callback(true)
						return
					}
					callback(false, ret)
				})
		}).
		Next(func(callback waterfall.Callback, args ...any) {
			ack := args[0].(*mymsg.SMQuerySceneAck)
			app.Request(ns, "logic.logicremote.reqswitchline", p.GetLogicId(), &mymsg.LogicReqSwitchLine{
				UId:         p.uid,
				SceneServer: ack.ServiceId,
				SceneId:     ack.SceneId,
				CfgId:       ack.CfgId,
				Token:       ack.Token,
			}, func(err error, ret any) {
				if err != nil {
					l.L.Errorf("reqswitchline failed, err: %v", err)
					callback(true)
					return
				}
				callback(false)
			})
		}).
		Final(func(err bool, args ...any) {
			p.sceneSwitching = false
		}).
		Do()
}
