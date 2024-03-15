package sysfactory

import (
	"mmo/common/factory"
)

var systemFactory = factory.NewStringFactory()

func GetFactory() *factory.StringFactory {
	return systemFactory
}
