package actions

import (
	"mlrs/b3"
	"mlrs/b3/core"
)

type Failure struct {
	core.Action
}

func (node *Failure) OnTick(tick *core.Tick) b3.Status {
	return b3.FAILURE
}
