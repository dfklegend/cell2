package events

import (
	"mmo/common/entity"
)

type EventBeHit struct {
	Src entity.IComponent
	Tar entity.IComponent
	Dmg int32
}
