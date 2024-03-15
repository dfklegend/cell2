package common

import (
	"github.com/dfklegend/cell2/utils/common"
)

// 有超时的状态

type StateWithTimeout struct {
	state   int
	timeout int64
}

func NewStateWithTimeout() *StateWithTimeout {
	return &StateWithTimeout{}
}

// SetState timeout 0 代表永不超时
func (l *StateWithTimeout) SetState(state int, timeout int64) {
	l.state = state
	if timeout > 0 {
		l.timeout = common.NowMs() + timeout
	} else {
		l.timeout = 0
	}
}

func (l *StateWithTimeout) GetState() int {
	return l.state
}

func (l *StateWithTimeout) IsTimeout() bool {
	return l.timeout > 0 && common.NowMs() >= l.timeout
}
