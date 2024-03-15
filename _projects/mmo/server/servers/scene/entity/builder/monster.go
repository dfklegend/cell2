package builder

import (
	"mmo/common/entity"
	"mmo/modules/fight/common"
	"mmo/modules/utils"
	define3 "mmo/servers/scene/define"
	ctrl2 "mmo/servers/scene/entity/ai/ctrl"
	"mmo/servers/scene/entity/attr/monster"
	components2 "mmo/servers/scene/entity/components"
	define2 "mmo/servers/scene/entity/define"
)

type MonsterInfo struct {
	CfgId string
	Name  string
	Side  int
	Level int
	Pos   define3.Pos
}

// BaseUnit
// StepMove
// Transform
// Skill
// AI
// MonsterComponent

func CreateMonsterEntity(world entity.IWorld, info *MonsterInfo) entity.IEntity {
	e := createNewEntity(world)

	baseUnit := e.AddComponent(define2.BaseUnit, components2.NewBaseUnit(define3.UnitMonster)).(*components2.BaseUnit)
	baseUnit.Name = info.Name

	e.AddComponent(define2.MoveComponent, components2.NewStepMoveComponent(false))

	tran := components2.NewTransform()
	tran.SetPos(utils.MakeGridPos(info.Pos))
	e.AddComponent(define2.Transform, tran)
	e.AddComponent(define2.Skill, components2.NewSkillComponent())

	ai := e.AddComponent(define2.AI, components2.NewAIComponent()).(*components2.AIComponent)
	ai.InitAICtrl(ctrl2.NewAICtrl(ctrl2.Normal))
	e.AddComponent(define2.MonsterComponent, components2.NewMonsterComponent(info.CfgId))

	e.Prepare()
	// 初始化
	char := baseUnit.GetChar()
	char.ChangeLevel(info.Level)

	char.SetIntBaseValue(common.Side, info.Side)
	char.SetIntBaseValue(common.EnergyMax, 100)
	char.SetIntBaseValue(common.PhysicPower, 10)
	char.SetIntBaseValue(common.PhysicArmor, 1)
	char.SetIntBaseValue(common.HPMax, 1000)
	char.SetBaseValue(common.AttackSpeed, 1)

	char.AddEquipGroup(monster.NewAttrInitor(info.CfgId))

	//char.SetEquip(0, "电刀")
	//char.SetEquip(1, "鬼刀")
	//char.SetEquip(2, "泰坦")

	skills := char.GetSkillTable()
	skills.SetNormalAttackSkill("普攻")
	skills.AddSkill("英勇打击", 1)

	char.Born()

	e.Start()
	world.AddEntity(e)
	return e
}
