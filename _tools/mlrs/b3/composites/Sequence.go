package composites

import (
	"mlrs/b3"
	. "mlrs/b3/core"
)

type Sequence struct {
	Composite
}

func (node *Sequence) OnTick(tick *Tick) b3.Status {
	//fmt.Println("tick Sequence :", node.GetTitle())
	for i := 0; i < node.GetChildCount(); i++ {
		var status = node.GetChild(i).Execute(tick)
		if status != b3.SUCCESS {
			return status
		}
	}
	return b3.SUCCESS
}
