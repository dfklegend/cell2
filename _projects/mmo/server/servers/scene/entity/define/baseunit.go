package define

import (
	common2 "mmo/modules/fight/common"
	"mmo/servers/scene/define"
)

type IBaseUnit interface {
	GetUnitType() define.UnitType
	GetChar() common2.ICharacter
	IsDead() bool
}
