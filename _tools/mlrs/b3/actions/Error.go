package actions

import (
	"mlrs/b3"
	"mlrs/b3/core"
)

type Error struct {
	core.Action
}

func (node *Error) OnTick(tick *core.Tick) b3.Status {
	return b3.ERROR
}
