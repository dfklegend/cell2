package center

import (
	"github.com/dfklegend/cell2/utils/common"
	l "github.com/dfklegend/cell2/utils/logger"
)

// PlayerTransactionLock
// 限制玩家上下线，切线等重要流程，防止重入
// 超时后，自动解锁，超时时间一般设成10分钟
// 超时解锁，输出重点错误log
type PlayerTransactionLock struct {
	lock bool
	// 超时时间
	timeout int64

	uid int64
	// lock/unlock 成对出现
	reason int
}

func newPlayerTransactionLock() *PlayerTransactionLock {
	return &PlayerTransactionLock{}
}

func (t *PlayerTransactionLock) Init(uid int64) {
	t.uid = uid
}

func (t *PlayerTransactionLock) IsLocked() bool {
	return t.lock && common.NowMs() > t.timeout
}

func (t *PlayerTransactionLock) Lock(reason int, timeout int64) bool {
	now := common.NowMs()
	// 已经锁定情况下，判断是否超时
	if t.lock {
		if now < t.timeout {
			return false
		}
		t.timeoutUnlock()
	}

	t.lock = true
	t.reason = reason
	t.timeout = now + timeout
	return true
}

func (t *PlayerTransactionLock) timeoutUnlock() {
	t.lock = false
	l.L.Errorf("PlayerTransactionLock timeout, uid: %v lock reason: %v", t.uid, t.reason)
}

func (t *PlayerTransactionLock) Unlock(reason int) bool {
	if !t.lock {
		l.L.Errorf("try to unlock a unlocked lock, uid: %v reason : %v ", t.uid, reason)
		return false
	}

	if reason != t.reason {
		l.L.Errorf("try to unlock mismatch reason, uid: %v reason cur: %v call: %v", t.uid, t.reason, reason)
		return false
	}
	t.lock = false
	return true
}
