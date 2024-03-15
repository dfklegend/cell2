package builder

import (
	"mmo/common/entity"
	"mmo/modules/fight/common"
	define3 "mmo/servers/scene/define"
	ctrl2 "mmo/servers/scene/entity/ai/ctrl"
	components2 "mmo/servers/scene/entity/components"
	define2 "mmo/servers/scene/entity/define"
)

func CreateCardMonsterEntity(world entity.IWorld, info *MonsterInfo, cb func(c common.ICharacter)) entity.IEntity {
	e := createNewEntity(world)

	baseUnit := e.AddComponent(define2.BaseUnit, components2.NewBaseUnit(define3.UnitMonster)).(*components2.BaseUnit)
	baseUnit.Name = info.Name

	tran := components2.NewTransform()
	tran.SetPos(info.Pos)
	e.AddComponent(define2.Transform, tran)
	e.AddComponent(define2.Skill, components2.NewSkillComponent())

	ai := e.AddComponent(define2.AI, components2.NewAIComponent()).(*components2.AIComponent)
	ai.InitAICtrl(ctrl2.NewAICtrl(ctrl2.CardNormal))

	e.Prepare()
	// 初始化
	char := baseUnit.GetChar()
	char.SetIntBaseValue(common.Side, info.Side)
	char.SetIntBaseValue(common.EnergyMax, 100)

	//char.ChangeLevel(100 + int(rand.Int31n(10)))
	//char.SetIntBaseValue(common.PhysicPower, 10)
	//char.SetIntBaseValue(common.PhysicArmor, 1)
	//char.AddEquipGroup(attrinitor.NewBaseModelInitor())
	//
	//skills := char.GetSkillTable()
	//skills.SetNormalAttackSkill("普攻")
	//skills.AddSkill("英勇打击", 1)
	//skills.AddSkill("嗜血", 1)
	//skills.AddSkill("致命打击", 1)
	if cb != nil {
		cb(char)
	}

	char.Born()

	e.Start()
	world.AddEntity(e)
	return e
}
