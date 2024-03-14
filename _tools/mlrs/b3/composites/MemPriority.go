package composites

import (
	"mlrs/b3"
	. "mlrs/b3/core"
)

type MemPriority struct {
	Composite
}

func (node *MemPriority) OnOpen(tick *Tick) {
	tick.Blackboard.Set("runningChild", 0, tick.GetTree().GetID(), node.GetID())
}

func (node *MemPriority) OnTick(tick *Tick) b3.Status {
	var child = tick.Blackboard.GetInt("runningChild", tick.GetTree().GetID(), node.GetID())
	for i := child; i < node.GetChildCount(); i++ {
		var status = node.GetChild(i).Execute(tick)

		if status != b3.FAILURE {
			if status == b3.RUNNING {
				tick.Blackboard.Set("runningChild", i, tick.GetTree().GetID(), node.GetID())
			}

			return status
		}
	}
	return b3.FAILURE
}
