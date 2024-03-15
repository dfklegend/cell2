package center

import (
	"mmo/common"
)

type Player struct {
	UId int64
	// 前端服务器
	FrontId string
	NetId   uint32

	logicId string

	state *common.StateWithTimeout
	lock  *PlayerTransactionLock

	// 客户端连接断开
	// 上线过程中收到断线消息后，需要等待
	// 上线完毕才能执行下线
	clientConnBroken bool

	FirstKickTime int64
	KickTimes     int
}

func NewPlayer(uid int64, frontId string, netId uint32) *Player {
	p := &Player{
		state:   common.NewStateWithTimeout(),
		lock:    newPlayerTransactionLock(),
		UId:     uid,
		FrontId: frontId,
		NetId:   netId,
	}
	p.SetState(Init, 0)
	p.lock.Init(uid)
	return p
}

func (p *Player) SetState(state int, timeout int64) {
	p.state.SetState(state, timeout)
}

// ReInit
// 重新设置数据
func (p *Player) ReInit(frontId string, netId uint32) {
	p.FrontId = frontId
	p.NetId = netId
	p.clientConnBroken = false
	p.KickTimes = 0
}

func (p *Player) ResetSession() {
	p.FrontId = ""
	p.NetId = 0
}

func (p *Player) IsSessionClosed() bool {
	return p.NetId == 0
}

func (p *Player) GetState() int {
	return p.state.GetState()
}

func (p *Player) IsTimeout() bool {
	return p.state.IsTimeout()
}

func (p *Player) NeedRemove() bool {
	return p.state.GetState() == WaitRemove
}

func (p *Player) TransactionLock(reason int, timeout int64) bool {
	return p.lock.Lock(reason, timeout)
}

func (p *Player) TransactionUnlock(reason int) bool {
	return p.lock.Unlock(reason)
}

func (p *Player) OnSessionClose() {
	p.ResetSession()
	p.clientConnBroken = true
}

func (p *Player) OnLogined(logicId string) {
	p.logicId = logicId
	p.SetState(Logined, 0)
}

func (p *Player) GetLogicId() string {
	return p.logicId
}
