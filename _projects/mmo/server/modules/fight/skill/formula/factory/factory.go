package factory

import (
	"mmo/common/factory"
)

var systemFactory = factory.NewIntFactory()

func GetFormulaFactory() *factory.IntFactory {
	return systemFactory
}
