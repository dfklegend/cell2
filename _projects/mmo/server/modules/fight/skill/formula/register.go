package formula

import (
	factory2 "mmo/common/factory"
	"mmo/modules/fight/skill/define"
	"mmo/modules/fight/skill/formula/factory"
)

func Register(t int, formula define.IFormula) {
	factory.GetFormulaFactory().RegisterFunc(t, func(args ...any) factory2.IObject {
		return formula
	})
}

func RegisterAll() {
	// register by init
}
