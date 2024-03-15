package effect

import (
	factory2 "mmo/common/factory"
	"mmo/modules/fight/skill/define"
	"mmo/modules/fight/skill/effect/factory"
)

func Register(t int, effect define.ISkillEffect) {
	factory.GetOpFactory().RegisterFunc(t, func(args ...any) factory2.IObject {
		return effect
	})
}

func RegisterAll() {
	// register by init
}
