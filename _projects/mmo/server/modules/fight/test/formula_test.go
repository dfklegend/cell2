package test

import (
	"testing"

	"github.com/dfklegend/cell2/utils/event/light"
	"github.com/stretchr/testify/assert"

	"mmo/modules/csv/entry"
	common2 "mmo/modules/fight/common"
	"mmo/modules/fight/common/skilldef"
	"mmo/modules/fight/skill"
	"mmo/modules/fight/skill/formula"
	"mmo/modules/fight/utils"
)

// 	测试基础的物理结算
// 100+ 0.5*100 - 30
func TestBasePhysic(t *testing.T) {
	cfg := &entry.Skill{}

	cfg.DmgType.Type = skilldef.DmgTypePhysic
	cfg.BaseDmg = 100
	cfg.AD = 0.5

	timeProvider := utils.NewTestTimeProvider()
	world := newTestWorld(timeProvider, nil)

	b1 := createCharBuilder()
	src := b1.Build()
	src.Init(0, world, light.NewEventCenter())
	src.SetIntBaseValue(common2.PhysicPower, 100)
	src.SetIntBaseValue(common2.PhysicArmor, 30)

	skill1 := skill.NewSkill("", 1, cfg, src)

	formula := formula.Normal{}

	dmg := skill.AllocDmgInstance()
	formula.Apply(skill1, src, src, dmg)

	assert.Equal(t, 120, dmg.Dmg)

	skill.FreeDmgInstance(dmg)
}

// 100+ 0.5*(100+10) - 30
func TestPlugin(t *testing.T) {
	cfg := &entry.Skill{}

	cfg.DmgType.Type = skilldef.DmgTypePhysic
	cfg.BaseDmg = 100
	cfg.AD = 0.5

	timeProvider := utils.NewTestTimeProvider()
	world := newTestWorld(timeProvider, nil)

	b1 := createCharBuilder()

	src := b1.Build()
	src.Init(0, world, light.NewEventCenter())
	src.SetIntBaseValue(common2.PhysicPower, 100)
	src.SetIntBaseValue(common2.PhysicArmor, 30)

	skill1 := skill.NewSkill("", 1, cfg, src)
	plugin1 := BonusPhysicPower{}
	plugin1.Start(src.GetEvents())

	formula := formula.Normal{}

	dmg := skill.AllocDmgInstance()
	formula.Apply(skill1, src, src, dmg)

	// 技能ad加成为50%
	assert.Equal(t, 125, dmg.Dmg)

	skill.FreeDmgInstance(dmg)
}
