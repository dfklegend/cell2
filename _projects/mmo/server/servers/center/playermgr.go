package center

import (
	"time"

	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/service"
	l "github.com/dfklegend/cell2/utils/logger"

	"mmo/common/define"
	mymsg "mmo/messages"
)

// 某些状态的超时时间

const (
	LoginTimeout  int64 = 2 * 60 * 1000
	LogoutTimeout int64 = 30 * 60 * 1000
)

/*
术语整理
login
    登录，客户端连接后，发起登录请求，验证通过可以进入游戏
offline
    客户端离线，服务器角色还存在
reonline
    offline后，客户端再次建立连接
logout
    角色登出，offline一段时间后，基于所在场景规则，角色会登出
    彻底离开游戏世界
    logout时，会存储角色数据
*/

type PlayerMgr struct {
	players     map[int64]*Player
	ns          *service.NodeService
	kickWaitMgr *KickWaitTaskMgr
}

func NewPlayerMgr(ns *service.NodeService) *PlayerMgr {
	m := &PlayerMgr{
		ns:      ns,
		players: make(map[int64]*Player),
	}
	m.kickWaitMgr = newKickWaitTaskMgr(ns, m)
	return m
}

func (m *PlayerMgr) Start() {
	m.ns.GetRunService().GetTimerMgr().AddTimer(1*time.Second, func(args ...interface{}) {
		m.update()
	})
}

func (m *PlayerMgr) update() {
	var toRemove []int64
	for k, v := range m.players {
		m.updatePlayer(k, v)

		if v.NeedRemove() {
			if toRemove == nil {
				toRemove = make([]int64, 0)
			}
			toRemove = append(toRemove, k)
		}
	}

	if toRemove != nil {
		for _, v := range toRemove {
			m.removePlayer(v)
		}
	}
}

// ReqLogin handler->
func (m *PlayerMgr) ReqLogin(uid int64, frontId string, netId uint32, kickPrev bool, cbFunc apientry.HandlerCBFunc) {
	m.ns.GetLogger().Infof("ReqLogin: %v, %v, %v", uid, frontId, netId)
	player := m.players[uid]
	if player != nil {
		if player.IsSessionClosed() {
			// 已经掉线，那么如果登录完毕，则走重连流程
			// 否则，返回失败(等待前一个登录流程完毕)
			if player.GetState() == Logined {

				// 前一个session close了, reconnect
				m.doReconnect(player, frontId, netId, cbFunc)
				return
			}

			m.ns.GetLogger().Infof("ReqLogin failed, already has a player %v in state: %v", uid, player.GetState())
			// 不是登录状态，暂时不能重连，等待前面事务完成
			apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.CenterReqLoginAck{
				Code: int32(define.ErrAlreadyOnline),
			})
			return
		}
		l.L.Infof("already has player: %v", uid)
		if kickPrev {
			// 注: 连续点击登录情况下，原连接客户端可能会收到多个kick消息
			// 主要原因是客户端延迟断开
			m.tryKick(player)
		}

		if player.GetState() == Logined {
			// 为了保证体验，可以等待被踢完，自动继续登录
			// 自动踢掉前一个kickWait
			m.kickWaitMgr.AddTask(uid, frontId, netId, cbFunc)
			return
		}

		// 如果是登录中，那就只能先返回错误
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.CenterReqLoginAck{
			Code: int32(define.ErrAlreadyOnline),
		})
		return
	}

	player = NewPlayer(uid, frontId, netId)
	player.SetState(Logining, LoginTimeout)

	m.players[uid] = player

	if !player.TransactionLock(TransactionLogin, TransactionLockLoginTimeout) {
		m.ns.GetLogger().Infof("TransactionLock(TransactionLogin) failed")
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.CenterReqLoginAck{
			Code: int32(define.ErrSystemBusy),
		})
		return
	}
	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.CenterReqLoginAck{
		Code: int32(define.Succ),
	})
}

func (m *PlayerMgr) doReconnect(player *Player, frontId string, netId uint32, cbFunc apientry.HandlerCBFunc) {
	// 角色不能处于锁定状态
	logger := m.ns.GetLogger()
	if !player.TransactionLock(TransactionReonline, TransactionLockTimeout) {
		logger.Infof("TransactionLock(TransactionReonline) failed!")
		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.CenterReqLoginAck{
			Code: int32(define.ErrSystemBusy),
		})
		return
	}

	player.ReInit(frontId, netId)
	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.CenterReqLoginAck{
		Code:        int32(define.Succ),
		IsReconnect: true,
		LogicId:     player.GetLogicId(),
	})
}

// tryKick
func (m *PlayerMgr) tryKick(player *Player) {
	l.L.Infof("kick player: %v, %v, %v", player.UId, player.FrontId, player.NetId)
	player.KickTimes++
	// 踢掉前端连接
	app.Kick(m.ns, player.FrontId, player.NetId, nil)
}

func (m *PlayerMgr) kickByNetId(uid int64, frontId string, netId uint32) {
	l.L.Infof("kickByNetId player: %v, %v,%v", uid, frontId, netId)
	// 踢掉前端连接
	app.Kick(m.ns, frontId, netId, nil)
}

// OnClientSessionClosed handler->
// 玩家连接断开
func (m *PlayerMgr) OnClientSessionClosed(uid int64) {
	l.L.Infof("OnClientSessionClosed, uid: %v", uid)

	player := m.players[uid]
	if player == nil {
		return
	}

	player.OnSessionClose()

	// 如果登录成功了
	if player.GetState() == Logined {
		m.sendOnOfflineAndCheckKickWait(uid, player.logicId)
	}

	// 否则等待登录成功了，再通知其离线
	if player.GetState() == Logining {
		l.L.Infof("OnClientSessionClosed when logining, uid: %v", uid)
	}
}

func (m *PlayerMgr) sendOnOffline(uid int64, logicId string, cb func(succ bool)) {
	app.Request(m.ns, "logic.logicremote.onoffline", logicId, &mymsg.OnOffline{
		UId: uid,
	}, func(err error, raw interface{}) {
		if err != nil {
			l.L.Errorf("logicremote.onoffline err: %v", err)
		}
		cb(err == nil)
	})
}

func (m *PlayerMgr) sendOnOfflineAndCheckKickWait(uid int64, logicId string) {
	// 简化，避免流程交叉: 等offline都处理完了，再触发kickWait
	m.sendOnOffline(uid, logicId, func(succ bool) {
		m.kickWaitMgr.OnSessionClose(uid)
	})
}

// OnLogicLogined logic->center
// 逻辑服
func (m *PlayerMgr) OnLogicLogined(uid int64, logicId string) {
	player := m.players[uid]
	if player == nil {
		return
	}
	if player.GetState() != Logining {
		l.L.Warnf("got logicLogined not in logining, state: %v", player.GetState())
	}

	player.OnLogined(logicId)
	player.TransactionUnlock(TransactionLogin)

	if player.IsSessionClosed() {
		l.L.Infof("client conn broken already when logined, uid: %v", uid)
		m.sendOnOfflineAndCheckKickWait(uid, logicId)
	}
}

// OnLogicReOnline
// logic -> center
func (m *PlayerMgr) OnLogicReOnline(uid int64) {
	player := m.players[uid]
	if player == nil {
		return
	}

	player.TransactionUnlock(TransactionReonline)
}

func (m *PlayerMgr) ReqLogout(uid int64) bool {
	player := m.players[uid]
	if player == nil {
		return false
	}
	if !player.TransactionLock(TransactionLogout, TransactionLockTimeout) {
		return false
	}
	player.SetState(Logouting, LogoutTimeout)
	return true
}

// OnLogicLogout logic->center
// 玩家在逻辑服登出完毕
func (m *PlayerMgr) OnLogicLogout(uid int64) {
	player := m.players[uid]
	if player == nil {
		return
	}

	if !player.IsSessionClosed() {
		m.tryKick(player)
	}

	player.TransactionUnlock(TransactionLogout)
	player.SetState(WaitRemove, 0)
}

// OnLogicAbnormalLogout
// logic发现scene失联调用，就是没有一个正确的scene对象时，会删除自己并自动logout
func (m *PlayerMgr) OnLogicAbnormalLogout(uid int64) {
	player := m.players[uid]
	if player == nil {
		return
	}

	if player.GetState() != Logined {
		m.ns.GetLogger().Warnf("%v OnLogicAbnormalLogout call when not in state logined", uid)
	}

	if !player.IsSessionClosed() {
		m.tryKick(player)
	}

	player.SetState(WaitRemove, 0)
}

func (m *PlayerMgr) removePlayer(uid int64) {
	m.ns.GetLogger().Infof("remove player: %v", uid)
	delete(m.players, uid)
}

// return true: player need to remove
func (m *PlayerMgr) updatePlayer(uid int64, p *Player) {
	if p.IsTimeout() {
		m.onPlayerTimeout(uid, p)
	}
}

func (m *PlayerMgr) onPlayerTimeout(uid int64, p *Player) {
	state := p.GetState()
	if state == Logining {
		m.onLoginningTimeout(uid, p)
		return
	}
	if state == Logouting {
		m.onLogoutingTimeout(uid, p)
		return
	}
}

func (m *PlayerMgr) onLoginningTimeout(uid int64, p *Player) {
	// 极端情况，某个时段连接不可用(比如logic crashed)
	// 没收到OnLogicLogined消息
	// 理论上 logic自己释放自己
	// 这里也是自己释放自己
	m.ns.GetLogger().Errorf("critical error onLoginningTimeout: %v", uid)
	p.SetState(WaitRemove, 0)
}

func (m *PlayerMgr) onLogoutingTimeout(uid int64, p *Player) {
	m.ns.GetLogger().Errorf("critical error onLogoutingTimeout: %v", uid)
	// 下线超时，更严重，可能说明数据没存储
	// 当前超时设置为30分钟，基本上只有可能是logic crashed
	p.SetState(WaitRemove, 0)
}

// OnGateKeepAlive
// gate定期同步其正常在线的player
func (m *PlayerMgr) OnGateKeepAlive() {

}

// 切线
func (m *PlayerMgr) ReqSwitchLine(uid int64) bool {
	player := m.players[uid]
	if player == nil {
		return false
	}

	if player.GetState() != Logined {
		m.ns.GetLogger().Errorf("switchline failed: %v, bad state: %v", uid, player.GetState())
		return false
	}

	if !player.TransactionLock(TransactionSwitchLine, TransactionLockTimeout) {
		m.ns.GetLogger().Errorf("switchline transactionLock failed: %v", uid)
		return false
	}
	player.SetState(SwitchLine, 0)
	return true
}

func (m *PlayerMgr) OnSwitchLineEnd(uid int64, succ bool) bool {
	player := m.players[uid]
	if player == nil {
		return false
	}

	if player.GetState() != SwitchLine {
		m.ns.GetLogger().Errorf("OnSwitchLineEnd failed: %v, bad state: %v", uid, player.GetState())
		return false
	}

	if !player.TransactionUnlock(TransactionSwitchLine) {
		m.ns.GetLogger().Errorf("switchline transactionLock failed: %v", uid)
		return false
	}

	player.SetState(Logined, 0)
	return true
}
