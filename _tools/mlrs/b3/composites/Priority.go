package composites

import (
	"mlrs/b3"
	. "mlrs/b3/core"
)

type Priority struct {
	Composite
}

func (node *Priority) OnTick(tick *Tick) b3.Status {
	for i := 0; i < node.GetChildCount(); i++ {
		var status = node.GetChild(i).Execute(tick)
		if status != b3.FAILURE {
			return status
		}
	}
	return b3.FAILURE
}
