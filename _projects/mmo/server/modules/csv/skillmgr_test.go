package csv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadSkill(t *testing.T) {
	cfg := NewSkillCfg()
	cfg.LoadFromFile("skill.csv")
	e1 := cfg.GetEntry("普攻")
	assert.Equal(t, true, e1 != nil)
}
