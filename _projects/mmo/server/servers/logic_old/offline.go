package logic_old

import (
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/service"
	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/waterfall"

	"mmo/common/define"
	mymsg "mmo/messages"
)

func ReqOffline(ns *service.NodeService, mgr *PlayerMgr, uid int64, cb func(err bool)) {
	// 如果在场景(战斗中)通知对方离线(战斗结果不要了)
	// . 请求存储
	// . 返回

	player := mgr.GetPlayer(uid)
	if player == nil {
		l.L.Warnf("can not find player: %v", uid)
		cb(false)
		return
	}

	// 存在风险: 下线过程中，如果有获取物品，将丢失
	player.BeginLogout()

	waterfall.NewBuilder(ns.GetRunService().GetScheduler()).
		Next(func(callback waterfall.Callback, args ...interface{}) {

			//callback := args[0].(waterfall.Callback)
			if player.IsInScene() {
				player.ReqLeaveCurScene(func(succ bool) {
					callback(false)
				})
			} else {
				callback(false)
			}

		}).
		Next(func(callback waterfall.Callback, args ...interface{}) {

			// saveplayer
			if !player.IsDirt() {
				callback(false)
				return
			}

			player.ClearDirt()
			l.L.Infof("saveplayer: %v", uid)
			app.Request(ns, "db.dbremote.saveplayer", nil, &mymsg.DBSavePlayer{
				Info: player.MakeData(),
			}, func(err error, raw interface{}) {
				if err != nil {
					l.L.Errorf("dbremote.saveplayer error: %v", err)
					callback(true)
					return
				}
				ack := raw.(*mymsg.NormalAck)
				if ack == nil || ack.Code != int32(define.Succ) {
					l.L.Errorf("dbremote.saveplayer error ")
					callback(true)
					return
				}

				// 后面考虑一种状态，会再次存储，存储之后自动offline
				if player.IsDirt() {
					l.L.Warnf("player %v dirt again!", player.UId)
				}
				callback(false)
			})

		}).
		Next(func(callback waterfall.Callback, args ...interface{}) {

			//l.L.Infof("centerremote.onlogicoffline: %v", uid)
			//app.Request(ns, "center.centerremote.onlogicoffline", nil,
			//	&mymsg.OnLogicOffline{
			//		UId: uid,
			//	}, func(err error, raw interface{}) {
			//		callback(false)
			//	})

		}).
		Next(func(callback waterfall.Callback, args ...interface{}) {

			//
			mgr.DestroyPlayer(uid)
			callback(false)

		}).
		Final(func(err bool, args ...interface{}) {
			cb(err)
		}).
		Do()
}
