package test

import (
	"log"
	"testing"

	"github.com/dfklegend/cell2/utils/event/light"
	"github.com/stretchr/testify/assert"

	"mmo/modules/csv"
	common2 "mmo/modules/fight/common"
	"mmo/modules/fight/skill/formula"
	"mmo/modules/fight/utils"
)

func TestEquip(t *testing.T) {
	formula.RegisterAll()

	csv.LoadAll("./testdata/csv")

	timeProvider := utils.NewTestTimeProvider()
	world := newTestWorld(timeProvider, nil)

	b1 := createCharBuilder()

	c1 := b1.Build()
	c1.Init(0, world, light.NewEventCenter())

	c1.GetAttr(common2.HPMax).SetBase(100)
	c1.Born()

	c1.SetEquip(0, "短剑_1")
	log.Println(c1.GetIntValue(common2.HPMax))
	assert.Equal(t, 150, c1.GetIntValue(common2.HPMax))

	c1.SetEquip(1, "盾牌_1")
	assert.Equal(t, 1870, c1.GetIntValue(common2.HPMax))

	c1.SetEquip(0, "盾牌_1")
	// 替换掉了短剑_1
	assert.Equal(t, 2940, c1.GetIntValue(common2.HPMax))

	c1.RemoveEquip(0)
	c1.RemoveEquip(1)

	assert.Equal(t, 100, c1.GetIntValue(common2.HPMax))
}
