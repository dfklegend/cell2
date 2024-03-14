package decorators

import (
	"mlrs/b3"
	. "mlrs/b3/config"
	. "mlrs/b3/core"
)

type Limiter struct {
	Decorator
	maxLoop int
}

func (node *Limiter) Initialize(setting *BTNodeCfg) {
	node.Decorator.Initialize(setting)
	node.maxLoop = setting.GetPropertyAsInt("maxLoop")
	if node.maxLoop < 1 {
		panic("maxLoop parameter in MaxTime decorator is an obligatory parameter")
	}
}

func (node *Limiter) OnTick(tick *Tick) b3.Status {
	if node.GetChild() == nil {
		return b3.ERROR
	}
	var i = tick.Blackboard.GetInt("i", tick.GetTree().GetID(), node.GetID())
	if i < node.maxLoop {
		var status = node.GetChild().Execute(tick)
		if status == b3.SUCCESS || status == b3.FAILURE {
			tick.Blackboard.Set("i", i+1, tick.GetTree().GetID(), node.GetID())
		}
		return status
	}

	return b3.FAILURE
}
