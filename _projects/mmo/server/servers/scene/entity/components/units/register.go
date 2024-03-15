package units

import (
	factory2 "mmo/common/factory"
	"mmo/servers/scene/define"
	"mmo/servers/scene/entity/components"
)

func Register() {
	factory := components.UnitImplFactory

	factory.RegisterFunc(int(define.UnitCamera), func(args ...any) factory2.IObject {
		return NewStaticUnitImpl()
	})

	factory.RegisterFunc(int(define.UnitExit), func(args ...any) factory2.IObject {
		return NewStaticUnitImpl()
	})

	factory.RegisterFunc(int(define.UnitMonster), func(args ...any) factory2.IObject {
		return NewMonsterUnitImpl()
	})

	factory.RegisterFunc(int(define.UnitAvatar), func(args ...any) factory2.IObject {
		return NewPlayerUnitImpl()
	})
}
