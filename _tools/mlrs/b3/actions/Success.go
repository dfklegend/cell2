package actions

import (
	"mlrs/b3"
	. "mlrs/b3/core"
)

type Success struct {
	Action
}

func (node *Success) OnTick(tick *Tick) b3.Status {
	return b3.SUCCESS
}
