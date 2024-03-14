package decorators

import (
	"mlrs/b3"
	. "mlrs/b3/config"
	. "mlrs/b3/core"
)

type RepeatUntilFailure struct {
	Decorator
	maxLoop int
}

func (node *RepeatUntilFailure) Initialize(setting *BTNodeCfg) {
	node.Decorator.Initialize(setting)
	node.maxLoop = setting.GetPropertyAsInt("maxLoop")
	if node.maxLoop < 1 {
		panic("maxLoop parameter in MaxTime decorator is an obligatory parameter")
	}
}

func (node *RepeatUntilFailure) OnOpen(tick *Tick) {
	tick.Blackboard.Set("i", 0, tick.GetTree().GetID(), node.GetID())
}

func (node *RepeatUntilFailure) OnTick(tick *Tick) b3.Status {
	if node.GetChild() == nil {
		return b3.ERROR
	}
	var i = tick.Blackboard.GetInt("i", tick.GetTree().GetID(), node.GetID())
	var status = b3.ERROR
	for node.maxLoop < 0 || i < node.maxLoop {
		status = node.GetChild().Execute(tick)
		if status == b3.SUCCESS {
			i++
		} else {
			break
		}
	}

	tick.Blackboard.Set("i", i, tick.GetTree().GetID(), node.GetID())
	return status
}
