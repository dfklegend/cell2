package decorators

import (
	"mlrs/b3"
	. "mlrs/b3/core"
)

type Inverter struct {
	Decorator
}

func (node *Inverter) OnTick(tick *Tick) b3.Status {
	if node.GetChild() == nil {
		return b3.ERROR
	}

	var status = node.GetChild().Execute(tick)
	if status == b3.SUCCESS {
		status = b3.FAILURE
	} else if status == b3.FAILURE {
		status = b3.SUCCESS
	}

	return status
}
