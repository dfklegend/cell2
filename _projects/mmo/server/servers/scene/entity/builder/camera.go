package builder

import (
	"mmo/common/entity"
	"mmo/servers/scene/define"
	components2 "mmo/servers/scene/entity/components"
	define2 "mmo/servers/scene/entity/define"
)

// BaseUnit
// Camera

func CreateCameraEntity(world entity.IWorld,
	owner entity.EntityID, frontId string, netId uint32) entity.IEntity {
	e := createNewEntity(world)

	e.AddComponent(define2.BaseUnit, components2.NewBaseUnit(define.UnitCamera))
	e.AddComponent(define2.Camera, components2.NewCamera(owner, frontId, netId))

	e.Prepare()
	e.Start()

	world.AddEntity(e)
	return e
}
