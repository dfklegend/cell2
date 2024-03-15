package define

import (
	define3 "mmo/servers/scene/define"
)

type IMoveComponent interface {
	MoveTo(tar define3.Pos)
	StopMove()
	IsMoving() bool
}
