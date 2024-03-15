package space

import (
	"mmo/common/entity"
	"mmo/servers/scene/define"
)

type EntityInfo struct {
	Id  entity.EntityID
	Pos define.Pos
}

func NewInfo() *EntityInfo {
	return &EntityInfo{}
}
