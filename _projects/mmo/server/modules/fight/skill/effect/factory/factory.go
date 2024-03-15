package factory

import (
	"mmo/common/factory"
)

var opFactory = factory.NewIntFactory()

func GetOpFactory() *factory.IntFactory {
	return opFactory
}
