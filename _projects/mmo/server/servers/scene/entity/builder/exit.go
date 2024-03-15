package builder

import (
	"mmo/common/entity"
	define3 "mmo/servers/scene/define"
	components2 "mmo/servers/scene/entity/components"
	define2 "mmo/servers/scene/entity/define"
)

type ExitInfo struct {
	Pos define3.Pos

	TarCfgId int32
	TarPos   define3.Pos
	Radius   float32
}

func CreateExitEntity(world entity.IWorld, info *ExitInfo) entity.IEntity {
	e := createNewEntity(world)

	baseUnit := e.AddComponent(define2.BaseUnit, components2.NewBaseUnit(define3.UnitExit)).(*components2.BaseUnit)
	baseUnit.Name = "场景切换点"

	tran := components2.NewTransform()
	tran.SetPos(info.Pos)
	e.AddComponent(define2.Transform, tran)
	e.AddComponent(define2.ExitComponent, components2.NewExitComponent(info.TarCfgId, info.TarPos, info.Radius))

	e.Prepare()
	e.Start()
	world.AddEntity(e)
	return e
}
