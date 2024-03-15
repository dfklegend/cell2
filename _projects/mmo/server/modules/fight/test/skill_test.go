package test

import (
	"testing"

	"github.com/dfklegend/cell2/utils/event/light"
	"github.com/stretchr/testify/assert"

	"mmo/modules/csv"
	common2 "mmo/modules/fight/common"
	"mmo/modules/fight/skill"
	"mmo/modules/fight/skill/formula"
	"mmo/modules/fight/utils"
)

//func

func TestBase(t *testing.T) {
	formula.RegisterAll()

	csv.LoadAll("./testdata/csv")
	cfg := csv.Skill.GetEntry("普攻")

	timeProvider := utils.NewTestTimeProvider()
	world := newTestWorld(timeProvider, nil)

	b1 := createCharBuilder()
	c1 := b1.Build()
	c1.Init(0, world, light.NewEventCenter())

	c1.GetAttr(common2.PhysicPower).SetBase(1)
	c1.GetAttr(common2.HPMax).SetBase(100)
	c1.Born()
	world.AddChar(c1)

	skill1 := skill.NewSkill("", 1, cfg, c1)
	skill1.SetTar(c1.GetId())
	skill1.SetWorld(world)

	skill1.Start()
	skill1.Update()

	timeProvider.SetNow(10000)

	skill1.Update()
	assert.Equal(t, false, skill1.IsOver())

	skill1.Update()
	assert.Equal(t, true, skill1.IsOver())
	assert.Equal(t, 99, c1.GetHP())
}

// 测试subskills
func TestSubSkills(t *testing.T) {
	formula.RegisterAll()

	csv.LoadAll("./testdata/csv")
	cfg := csv.Skill.GetEntry("菲奥娜技能")

	timeProvider := utils.NewTestTimeProvider()
	world := newTestWorld(timeProvider, nil)

	b1 := createCharBuilder()
	c1 := b1.Build()
	c1.Init(0, world, light.NewEventCenter())

	c1.GetAttr(common2.PhysicPower).SetBase(1)
	c1.GetAttr(common2.HPMax).SetBase(100)
	c1.Born()
	world.AddChar(c1)

	skill1 := skill.NewSkill("菲奥娜技能", 1, cfg, c1)
	skill1.SetTar(c1.GetId())
	skill1.SetWorld(world)

	skill1.Start()
	skill1.Update()

	timeProvider.SetNow(10000)
	skill1.Update()
	skill1.Update()
	skill1.Update()
	assert.Equal(t, true, skill1.IsOver())
}

// 测试subskills break
func TestSubSkillBreak(t *testing.T) {
	formula.RegisterAll()

	csv.LoadAll("./testdata/csv")
	cfg := csv.Skill.GetEntry("菲奥娜技能")

	timeProvider := utils.NewTestTimeProvider()
	world := newTestWorld(timeProvider, nil)

	b1 := createCharBuilder()
	c1 := b1.Build()
	c1.Init(0, world, light.NewEventCenter())

	c1.GetAttr(common2.PhysicPower).SetBase(1)
	c1.GetAttr(common2.HPMax).SetBase(100)
	c1.Born()
	world.AddChar(c1)

	skill1 := skill.NewSkill("菲奥娜技能", 1, cfg, c1)
	skill1.SetTar(c1.GetId())
	skill1.SetWorld(world)

	skill1.Start()
	skill1.Update()

	timeProvider.SetNow(10000)
	skill1.TestBreakSubSkill()
	assert.Equal(t, false, skill1.IsFailed())
	skill1.Update()
	assert.Equal(t, true, skill1.IsFailed())
	assert.Equal(t, skill.ReasonSubFailed, skill1.GetFailedReason())
}
