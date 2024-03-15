package scenecfg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	cfg := LoadSceneCfg("./", "1")
	assert.Equal(t, true, cfg != nil)
	assert.Equal(t, "1", cfg.Monsters[0].Id)
}

func Test2(t *testing.T) {
	cfg := LoadSceneCfg("./", "2")
	assert.Equal(t, true, cfg != nil)
	assert.Equal(t, true, cfg.Monsters == nil)
}
