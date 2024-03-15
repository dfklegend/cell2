package skill

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mmo/modules/fight/utils"

	"mmo/modules/csv"
)

func TestTable(t *testing.T) {
	csv.LoadAll("./testdata/csv")

	tb := NewSkillTable()
	tb.Init(nil, utils.NewTestTimeProvider())

	tb.AddSkill("skill1", 1)
	tb.AddSkill("skill1", 2)
	assert.Equal(t, 0, tb.findIndex("skill1"))
	tb.RemoveSkill("skill1")
	assert.Equal(t, -1, tb.findIndex("skill1"))
}
