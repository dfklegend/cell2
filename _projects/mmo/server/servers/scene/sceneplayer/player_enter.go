package sceneplayer

import (
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/utils/common"
	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/waterfall"

	define2 "mmo/common/define"
	mymsg "mmo/messages"
	"mmo/servers/scene/define"
)

func (p *ScenePlayer) EnterWorld(newPlayer bool, info *mymsg.PlayerInfo, switchLine bool) {
	if newPlayer {
		p.InitNewPlayerInfo()
	} else {
		p.LoadInfo(info)
	}

	p.OnEnterWorld(switchLine)
}

// OnEnterWorld
// switchLine: 是否切线进入(切线进入)
func (p *ScenePlayer) OnEnterWorld(switchLine bool) {
	p.ns.GetLogger().Infof("%v onEnterWorld switchLine: %v", p.uid, switchLine)
	p.nextKeepAlive = common.NowMs() + define2.SceneKeepAliveMs
	p.systemsOnEnterWorld(switchLine)
	p.events.Publish("onEnterWorld", switchLine)
}

func (p *ScenePlayer) ReEnter(frontId string, netId uint32) {
	// 断线重连
	p.frontId = frontId
	p.netId = netId
	p.OnReOnline()
}

func (p *ScenePlayer) Enter(switchLine bool, cb func(succ bool)) {
	p.ChangeState(define.LoadingInfo)

	ns := p.ns
	uid := p.GetId()
	waterfall.NewBuilder(ns.GetRunService().GetScheduler()).
		Next(func(callback waterfall.Callback, args ...interface{}) {
			app.Request(ns, "db.dbremote.loadplayer", nil, &mymsg.DBLoadPlayer{
				UId: uid,
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
				callback(false, ack)
			})
		}).
		Next(func(callback waterfall.Callback, args ...any) {
			// load info
			ack := args[0].(*mymsg.DBLoadPlayerAck)
			p.EnterWorld(ack.NewPlayer, ack.Info, switchLine)
			callback(false)
		}).
		Final(func(err bool, args ...any) {
			if err {
				p.ChangeState(define.WaitMgrDelete)
			}
			cb(!err)
		}).
		Do()
}
