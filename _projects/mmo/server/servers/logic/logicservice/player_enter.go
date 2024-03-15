package logicservice

import (
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/node/app"
	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/waterfall"

	"mmo/common/define"
	mymsg "mmo/messages"
)

// Enter
// 进入流程
// 进入，分配场景
// @cbFunc: 回调函数
//		参数(err error,result *mymsg.NormalAck)
func (p *LogicPlayer) Enter(msg *mymsg.LogicLoadPlayer, cbFunc apientry.HandlerCBFunc) {
	// . 请求db载入角色
	// . player.LoadInfo
	// . 请求进入场景
	// . 通知center上线

	l.L.Errorf("logicPlayer.EnterWorld")
	ns := p.ns
	waterfall.NewBuilder(ns.GetRunService().GetScheduler()).
		Next(func(callback waterfall.Callback, args ...interface{}) {
			// loadplayer 只是为了获取场景信息
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
				callback(false)
			})
		}).
		Next(func(callback waterfall.Callback, args ...interface{}) {
			// 请求进入场景
			// mmo规则，是尽量回到上一次场景，当前先创建新的，后续再处理
			p.EnterWorld()
			p.AllocEnterScene(int(define.ScenePlayerHome), define.Pos{}, false, false, func(succ bool) {
				if succ {
					callback(false)
				} else {
					callback(true)
				}
			})

		}).
		Next(func(callback waterfall.Callback, args ...interface{}) {
			// notify handler
			app.Request(ns, "center.centerremote.onlogiclogined", nil, &mymsg.OnLogicLogined{
				UId:     msg.UId,
				LogicId: ns.Name,
			}, func(err error, raw interface{}) {
				if err != nil {
					l.L.Errorf("centerremote.onlogiclogined error: %v", err)
					callback(true)
					return
				}

				ack := raw.(*mymsg.NormalAck)
				if ack == nil || ack.Code != int32(define.Succ) {
					l.L.Errorf("centerremote.onlogiclogined failed ")
					callback(true)
					return
				}
				callback(false)
			})
		}).
		Final(func(err bool, args ...interface{}) {
			// TODO: 如果失败，需要处理player对象
			if err {
				apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
					Code:  int32(define.ErrWithStr),
					Error: "logic login failed",
				})
			} else {
				apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
					Code: int32(define.Succ),
				})
			}
		}).
		Do()
}

func (p *LogicPlayer) OnOffline() {
	p.ns.GetLogger().Infof("%v offline", p.uid)
}

func (p *LogicPlayer) OnReOnline() {
	p.ns.GetLogger().Infof("%v reonline", p.uid)
}

// ReEnter
// 玩家断线重连
// @cbFunc: 回调函数
//		参数(err error,result *mymsg.NormalAck)
func (p *LogicPlayer) ReEnter(cbFunc apientry.HandlerCBFunc) {
	// . 通知scene reenter
	//logger := p.ns.GetLogger()

	ns := p.ns
	waterfall.NewBuilder(ns.GetRunService().GetScheduler()).
		Next(func(callback waterfall.Callback, args ...interface{}) {
			p.ReEnterScene(func(succ bool) {
				if succ {
					callback(false)
				} else {
					callback(true)
				}
			})
		}).
		Next(func(callback waterfall.Callback, args ...interface{}) {
			app.Request(ns, "center.centerremote.onlogicreonline", nil, &mymsg.OnLogicLogined{
				UId:     p.GetUId(),
				LogicId: ns.Name,
			}, func(err error, raw interface{}) {
				if err != nil {
					l.L.Errorf("centerremote.onlogicreonline error: %v", err)
					callback(true)
					return
				}

				ack := raw.(*mymsg.NormalAck)
				if ack == nil || ack.Code != int32(define.Succ) {
					l.L.Errorf("centerremote.onlogicreonline failed ")
					callback(true)
					return
				}
				callback(false)
			})
		}).
		Final(func(err bool, args ...interface{}) {
			if err {
				apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
					Code:  int32(define.ErrWithStr),
					Error: "logic reenter failed",
				})
			} else {
				apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
					Code: int32(define.Succ),
				})
			}
		}).
		Do()

}
