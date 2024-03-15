package builder

import (
	"math/rand"

	"mmo/common/entity"
	"mmo/modules/fight/common"
	"mmo/modules/utils"
	define3 "mmo/servers/scene/define"
	ctrl2 "mmo/servers/scene/entity/ai/ctrl"
	"mmo/servers/scene/entity/attr/player"
	components2 "mmo/servers/scene/entity/components"
	define2 "mmo/servers/scene/entity/define"
	"mmo/servers/scene/sceneplayer/systems/interfaces"
)

type PlayerInfo struct {
	Player define3.IPlayer
}

// CreateAvatarEntity 创建玩家Avatar对象
//
// 包含Components
// 		BaseUnit
// 		Move
// 		Transform
// 		Control
// 		Player
//		AI
// @param world
// @param info 创建所需信息
// @return entity 返回创建的entity
func CreateAvatarEntity(world entity.IWorld, info *PlayerInfo) entity.IEntity {
	e := createNewEntity(world)

	baseUnit := e.AddComponent(define2.BaseUnit, components2.NewBaseUnit(define3.UnitAvatar)).(*components2.BaseUnit)
	baseUnit.Name = "玩家"

	e.AddComponent(define2.MoveComponent, components2.NewStepMoveComponent(true))

	tran := components2.NewTransform()
	tran.SetPos(utils.MakeGridPos(define3.Pos{
		X: 10*rand.Float32() - 5,
		Z: 10*rand.Float32() - 5,
	}))
	e.AddComponent(define2.Transform, tran)
	e.AddComponent(define2.Skill, components2.NewSkillComponent())
	e.AddComponent(define2.Control, components2.NewControlComponent())
	e.AddComponent(define2.Player, components2.NewPlayerComponent(info.Player))

	ai := e.AddComponent(define2.AI, components2.NewAIComponent()).(*components2.AIComponent)
	ai.InitAICtrl(ctrl2.NewAICtrl(ctrl2.Player))

	e.Prepare()

	char := baseUnit.GetChar()
	char.ChangeLevel(int(info.Player.GetSystem(interfaces.BaseInfo).(interfaces.IBaseInfo).GetLevel()))
	char.SetIntBaseValue(common.Side, 1)
	char.SetIntBaseValue(common.EnergyMax, 100)
	char.SetIntBaseValue(common.PhysicPower, 30)
	char.SetIntBaseValue(common.PhysicArmor, 6)
	char.SetIntBaseValue(common.HPMax, 10000)
	char.SetBaseValue(common.AttackSpeed, 1)

	char.AddEquipGroup(player.NewAttrInitor())

	char.SetEquip(1, "鬼刀")

	skills := char.GetSkillTable()
	skills.SetNormalAttackSkill("普攻")
	skills.AddSkill("英勇打击", 1)

	char.Born()

	e.Start()

	world.AddEntity(e)
	return e
}
