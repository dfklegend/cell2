package logic_old

import (
	"time"

	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/service"
	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/timer"
)

const (
	Init = iota
	Running
	Retiring  // 退休中
	Retiring1 // 等待10s，如果还有在线玩家，一定是严重bug
	Retired   // 退休完毕
)

const (
	IntervalNotify      = 15
	IntervalTotalRetire = 5 * 60

	IntervalRetiring1Wait = 10
)

// RetireMgr
// logic的退休管理器
// 策略:
//		设置状态为退休中
// 		先踢掉在线空闲的玩家，不接受新的玩家进入
//		部分无法及时下线的玩家，比如: 战斗中
//		每30s关播，通知玩家
//		延迟5分钟之后，强制踢掉

type RetireMgr struct {
	state int
	mgr   *PlayerMgr

	ns          *service.NodeService
	timerUpdate timer.IdType

	tickNotify int
	tickTotal  int
}

func NewRetireMgr() *RetireMgr {
	return &RetireMgr{
		state: Init,
	}
}

func (r *RetireMgr) setState(state int) {
	r.state = state
}

func (r *RetireMgr) IsState(state int) bool {
	return r.state == state
}

func (r *RetireMgr) IsRunning() bool {
	return r.IsState(Running)
}

func (r *RetireMgr) Start(ns *service.NodeService, mgr *PlayerMgr) {
	r.ns = ns
	r.mgr = mgr
	r.setState(Running)
}

// Handler ICtrlCmdListener impl
func (r *RetireMgr) Handler(cmd string) string {
	if "queryretire" == cmd {
		return "ok"
	}

	if "retire" == cmd {
		r.beginRetire()
		return "ok"
	}
	return ""
}

func (r *RetireMgr) beginRetire() bool {
	if !r.IsState(Running) {
		l.L.Errorf("not in running can not retire")
		return false
	}

	r.setState(Retiring)
	l.L.Infof("%v -> Retiring", r.ns.Name)

	r.tryKickIdlePlayers()

	// start timers
	r.timerUpdate = r.ns.GetRunService().GetTimerMgr().AddTimer(time.Second, func(args ...interface{}) {
		r.update()
	})
	return true
}

func (r *RetireMgr) NotifyAllPlayers(content string) {
	r.mgr.Filter(func(p *Player) {
		p.PushBattleLog(content)
	})
}

func (r *RetireMgr) tryKickIdlePlayers() {
	r.mgr.Filter(func(p *Player) {
		if p.IsIdle() {
			r.mgr.Kick(p.UId)
		}
	})
}

func (r *RetireMgr) tryKickAllPlayers() {
	r.mgr.Filter(func(p *Player) {
		r.mgr.Kick(p.UId)
	})
}

func (r *RetireMgr) getActivePlayerNum() int {
	return r.mgr.GetPlayerNum()
}

func (r *RetireMgr) update() {
	if r.IsState(Retiring) {
		r.updateRetire()
		return
	}
	if r.IsState(Retiring1) {
		r.updateRetire1()
		return
	}
}

func (r *RetireMgr) canGoRetire1() bool {
	if r.tickTotal > 10 && r.getActivePlayerNum() == 0 {
		return true
	}
	if r.tickTotal > IntervalTotalRetire {
		return true
	}
	return false
}

func (r *RetireMgr) updateRetire() {
	r.tickNotify++
	r.tickTotal++

	if r.tickNotify > IntervalNotify {
		r.NotifyAllPlayers("服务器即将关闭，请尽快下线保存")
		r.tryKickIdlePlayers()
		r.tickNotify = 0
	}

	if r.canGoRetire1() {
		r.tickTotal = 0
		// 所有玩家都踢掉
		r.tryKickAllPlayers()
		r.setState(Retiring1)
		l.L.Infof("%v -> Retiring1", r.ns.Name)
	}
}

func (r *RetireMgr) updateRetire1() {
	r.tickTotal++

	if r.tickTotal > IntervalRetiring1Wait {
		// 已经等待10s了，应该差不多了
		if r.getActivePlayerNum() > 0 {
			l.L.Errorf("retireMgr still have players: %v", r.getActivePlayerNum())
		}

		r.setState(Retired)
		l.L.Infof("%v -> Retired", r.ns.Name)
		r.notifyServiceRetired()

		r.ns.GetRunService().GetTimerMgr().Cancel(r.timerUpdate)
		r.timerUpdate = 0
	}
}

func (r *RetireMgr) notifyServiceRetired() {
	app.NotifyServiceRetired(r.ns)
}
