package test

import (
	"testing"

	"github.com/dfklegend/cell2/utils/event/light"
	"github.com/stretchr/testify/assert"

	"mmo/modules/csv"
	common2 "mmo/modules/fight/common"
	"mmo/modules/fight/skill/formula"
	"mmo/modules/fight/utils"
)

func TestBufAddRemove(t *testing.T) {
	formula.RegisterAll()

	csv.LoadAll("./testdata/csv")

	timeProvider := utils.NewTestTimeProvider()
	world := newTestWorld(timeProvider, nil)

	b1 := createCharBuilder()

	c1 := b1.Build()
	c1.Init(0, world, light.NewEventCenter())

	c1.GetAttr(common2.HPMax).SetBase(100)
	c1.Born()

	c1.AddBuf(c1, c1, "英勇打击_减防", 1, 1)
	assert.Equal(t, -30, c1.GetAttr(common2.PhysicArmor).GetIntValue())
	c1.RemoveBuf("英勇打击_减防")
	assert.Equal(t, 0, c1.GetAttr(common2.PhysicArmor).GetIntValue())
}

func TestBufSpecialStatus(t *testing.T) {
	formula.RegisterAll()

	csv.LoadAll("./testdata/csv")

	timeProvider := utils.NewTestTimeProvider()
	world := newTestWorld(timeProvider, nil)

	b1 := createCharBuilder()

	c1 := b1.Build()
	c1.Init(0, world, light.NewEventCenter())

	c1.GetAttr(common2.HPMax).SetBase(100)
	c1.Born()

	c1.AddBuf(c1, c1, "菲奥娜_无敌", 1, 1)
	assert.Equal(t, true, c1.HasSpecialStatus(common2.SSInvincible))
	c1.RemoveBuf("菲奥娜_无敌")
	assert.Equal(t, false, c1.HasSpecialStatus(common2.SSInvincible))
}

func TestBufRemove_SameId(t *testing.T) {
	formula.RegisterAll()

	csv.LoadAll("./testdata/csv")

	timeProvider := utils.NewTestTimeProvider()
	world := newTestWorld(timeProvider, nil)

	b1 := createCharBuilder()
	c1 := b1.Build()
	c1.Init(0, world, light.NewEventCenter())

	c1.GetAttr(common2.HPMax).SetBase(100)
	c1.Born()

	c1.AddBuf(c1, c1, "英勇打击_减防", 1, 1)
	c1.RemoveBuf("英勇打击_减防")
	assert.Equal(t, 0, c1.GetAttr(common2.PhysicArmor).GetIntValue())

	c1.AddBuf(c1, c1, "英勇打击_减防", 1, 1)
	// 同id buf,多次删除(中间状态)
	// 此时存在一个待删除的和新加的buf
	c1.Update()
	assert.Equal(t, -30, c1.GetAttr(common2.PhysicArmor).GetIntValue())
	timeProvider.SetNow(10000)
	c1.Update()
	assert.Equal(t, 0, c1.GetAttr(common2.PhysicArmor).GetIntValue())
}

func TestBuf_RefreshNormal(t *testing.T) {
	formula.RegisterAll()

	csv.LoadAll("./testdata/csv")

	timeProvider := utils.NewTestTimeProvider()
	world := newTestWorld(timeProvider, nil)

	b1 := createCharBuilder()

	c1 := b1.Build()
	c1.Init(0, world, light.NewEventCenter())

	c1.GetAttr(common2.HPMax).SetBase(100)
	c1.Born()

	c1.AddBuf(c1, c1, "英勇打击_减防", 1, 1)
	c1.AddBuf(c1, c1, "英勇打击_减防", 1, 1)
	assert.Equal(t, -60, c1.GetAttr(common2.PhysicArmor).GetIntValue())

	c1.AddBuf(c1, c1, "英勇打击_减防", 1, 4)
	// 最多5层
	assert.Equal(t, -150, c1.GetAttr(common2.PhysicArmor).GetIntValue())

	c1.AddBuf(c1, c1, "英勇打击_减防", 1, 3)
	assert.Equal(t, -150, c1.GetAttr(common2.PhysicArmor).GetIntValue())
}

func TestBuf_Level(t *testing.T) {
	formula.RegisterAll()

	csv.LoadAll("./testdata/csv")

	timeProvider := utils.NewTestTimeProvider()
	world := newTestWorld(timeProvider, nil)

	b1 := createCharBuilder()

	c1 := b1.Build()
	c1.Init(0, world, light.NewEventCenter())

	c1.GetAttr(common2.HPMax).SetBase(100)
	c1.Born()

	c1.AddBuf(c1, c1, "英勇打击_减防", 2, 1)
	assert.Equal(t, -30, c1.GetAttr(common2.PhysicArmor).GetIntValue())

	c1.AddBuf(c1, c1, "英勇打击_减防", 2, 1)
	assert.Equal(t, -60, c1.GetAttr(common2.PhysicArmor).GetIntValue())

}

func TestBuf_RefreshLevel(t *testing.T) {
	formula.RegisterAll()

	csv.LoadAll("./testdata/csv")

	timeProvider := utils.NewTestTimeProvider()
	world := newTestWorld(timeProvider, nil)

	b1 := createCharBuilder()

	c1 := b1.Build()
	c1.Init(0, world, light.NewEventCenter())

	c1.GetAttr(common2.HPMax).SetBase(100)
	c1.Born()

	c1.AddBuf(c1, c1, "英勇打击_减防", 2, 2)

	// refresh low level will skipped
	c1.AddBuf(c1, c1, "英勇打击_减防", 1, 5)
	assert.Equal(t, -60, c1.GetAttr(common2.PhysicArmor).GetIntValue())

	c1.RemoveBuf("英勇打击_减防")
	c1.AddBuf(c1, c1, "英勇打击_减防", 1, 5)
	assert.Equal(t, -150, c1.GetAttr(common2.PhysicArmor).GetIntValue())

	// 高等级, 将挤掉低等级
	c1.AddBuf(c1, c1, "英勇打击_减防", 2, 1)
	assert.Equal(t, -30, c1.GetAttr(common2.PhysicArmor).GetIntValue())
}

func TestBuf_Get(t *testing.T) {
	formula.RegisterAll()

	csv.LoadAll("./testdata/csv")

	csv.Buf.GetEntry("盖伦技能_周期伤害")
}
