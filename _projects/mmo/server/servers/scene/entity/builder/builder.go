package builder

import (
	"math/rand"

	"mmo/common/entity"
	"mmo/common/entity/impl"
	define3 "mmo/servers/scene/define"
	components2 "mmo/servers/scene/entity/components"
	define2 "mmo/servers/scene/entity/define"
)

func createNewEntity(world entity.IWorld) entity.IEntity {
	id := world.AllocId()
	e := impl.NewEntity()
	e.SetWorld(world)
	e.SetId(id)
	return e
}

func CreateTestEntity(world entity.IWorld, info any) entity.IEntity {
	e := createNewEntity(world)

	e.AddComponent(define2.BaseUnit, components2.NewBaseUnit(define3.UnitTest))
	e.AddComponent(define2.MoveComponent, components2.NewStepMoveComponent(false))

	tran := components2.NewTransform()
	tran.SetPos(define3.Pos{
		X: (rand.Float32() - 0.5) * 2 * define3.MaxWidth,
		Z: (rand.Float32() - 0.5) * 2 * define3.MaxWidth,
	})
	e.AddComponent(define2.Transform, tran)

	e.Prepare()
	e.Start()

	world.AddEntity(e)
	return e
}
