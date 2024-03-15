package ai

import (
	fsm2 "mmo/servers/scene/entity/fsm"
)

var (
	aiStateFactory *fsm2.StateFactory = fsm2.NewStateFactory()
)

func GetStateFactory() fsm2.IFactory {
	return aiStateFactory
}
