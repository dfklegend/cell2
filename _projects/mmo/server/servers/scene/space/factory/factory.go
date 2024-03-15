package factory

import (
	"mmo/servers/scene/define"
)

type FuncCreateSpace func(args ...any) define.ISpace

type creatorRegistry struct {
	normal FuncCreateSpace
	zone   FuncCreateSpace
}

var registry = &creatorRegistry{}

func SetNormalCreator(f FuncCreateSpace) {
	registry.normal = f
}

func SetZoneCreator(f FuncCreateSpace) {
	registry.zone = f
}

func CreateNormalSpace(args ...any) define.ISpace {
	return registry.normal(args...)
}

func CreateZoneSpace(args ...any) define.ISpace {
	return registry.zone(args...)
}
