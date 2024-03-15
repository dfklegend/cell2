package sceneplayer

import (
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/utils/waterfall"

	define1 "mmo/common/define"
	mymsg "mmo/messages"
	"mmo/servers/scene/define"
)

// Leave
// 切线功能调用
func (p *ScenePlayer) Leave(cb func(succ bool)) {
	ns := p.ns
	logger := ns.GetLogger()

	if p.GetState() != define.Normal {
		logger.Errorf("player.leave error, bad state: %v", p.GetState())
		cb(false)
		return
	}

	logger.Infof("begin player.leave: %v", p.uid)

	// . 存储player
	// . 设置玩家可删除
	p.ChangeState(define.Logouting)

	// 离开当前场景
	p.LeaveCurScene()
	waterfall.NewBuilder(ns.GetRunService().GetScheduler()).
		Next(func(callback waterfall.Callback, args ...any) {
			p.savePlayer(func(succ bool) {
				callback(!succ)
			})
		}).
		Final(func(err bool, args ...any) {
			if err {
				logger.Errorf("critical error, player.leave %v failed!", p.uid)
				p.ChangeState(define.Normal)
				cb(false)
				return
			}
			p.ChangeState(define.WaitMgrDelete)
			cb(true)
		}).
		Do()
}

func (p *ScenePlayer) savePlayer(cb func(succ bool)) {
	if !p.IsDirt() {
		cb(true)
		return
	}

	ns := p.ns
	logger := ns.GetLogger()
	p.ClearDirt()
	logger.Infof("saveplayer: %v", p.uid)
	app.Request(ns, "db.dbremote.saveplayer", nil, &mymsg.DBSavePlayer{
		Info: p.MakeSaveInfo(),
	}, func(err error, raw interface{}) {
		if err != nil {
			logger.Errorf("dbremote.saveplayer error: %v", err)
			cb(false)
			return
		}
		ack := raw.(*mymsg.NormalAck)
		if ack == nil || ack.Code != int32(define1.Succ) {
			logger.Errorf("critical error dbremote.saveplayer error ")
			cb(false)
			return
		}

		// 避免此种情况(玩家下线或者leave前)
		// 不然会数据丢失
		if p.IsDirt() {
			logger.Errorf("critical error player %v dirt again when saving! will lost some player data", p.uid)
		}
		cb(true)
	})
}
