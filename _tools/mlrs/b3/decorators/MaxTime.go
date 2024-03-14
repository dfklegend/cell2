package decorators

import (
	"time"

	"mlrs/b3"
	. "mlrs/b3/config"
	. "mlrs/b3/core"
)

type MaxTime struct {
	Decorator
	maxTime int64
}

func (node *MaxTime) Initialize(setting *BTNodeCfg) {
	node.Decorator.Initialize(setting)
	node.maxTime = setting.GetPropertyAsInt64("maxTime")
	if node.maxTime < 1 {
		panic("maxTime parameter in Limiter decorator is an obligatory parameter")
	}
}

func (node *MaxTime) OnOpen(tick *Tick) {
	var startTime int64 = time.Now().UnixNano() / 1000000
	tick.Blackboard.Set("startTime", startTime, tick.GetTree().GetID(), node.GetID())
}

func (node *MaxTime) OnTick(tick *Tick) b3.Status {
	if node.GetChild() == nil {
		return b3.ERROR
	}
	var currTime int64 = time.Now().UnixNano() / 1000000
	var startTime int64 = tick.Blackboard.GetInt64("startTime", tick.GetTree().GetID(), node.GetID())
	var status = node.GetChild().Execute(tick)
	if currTime-startTime > node.maxTime {
		return b3.FAILURE
	}

	return status
}
