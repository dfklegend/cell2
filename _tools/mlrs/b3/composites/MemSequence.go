package composites

import (
	"mlrs/b3"
	. "mlrs/b3/core"
)

type MemSequence struct {
	Composite
}

func (node *MemSequence) OnOpen(tick *Tick) {
	tick.Blackboard.Set("runningChild", 0, tick.GetTree().GetID(), node.GetID())
}

func (node *MemSequence) OnTick(tick *Tick) b3.Status {
	var child = tick.Blackboard.GetInt("runningChild", tick.GetTree().GetID(), node.GetID())
	for i := child; i < node.GetChildCount(); i++ {
		var status = node.GetChild(i).Execute(tick)

		if status != b3.SUCCESS {
			if status == b3.RUNNING {
				tick.Blackboard.Set("runningChild", i, tick.GetTree().GetID(), node.GetID())
			}

			return status
		}
	}
	return b3.SUCCESS
}
