package actions

import (
	"mlrs/b3"
	"mlrs/b3/core"
)

type Running struct {
	core.Action
}

func (node *Running) OnTick(tick *core.Tick) b3.Status {
	return b3.FAILURE
}
