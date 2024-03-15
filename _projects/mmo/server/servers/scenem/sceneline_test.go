package scenem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FindIdle(t *testing.T) {
	lines := NewLines()

	line := NewSceneLine(0, 1, 2)
	lines.add(line)

	assert.Equal(t, int32(0), lines.FineIdleLineId())

	lines = NewLines()

	line = NewSceneLine(0, 1, 10)
	lines.add(line)

	line = NewSceneLine(0, 1, 0)
	lines.add(line)
	assert.Equal(t, int32(1), lines.FineIdleLineId())
}

func Test_FindIndex(t *testing.T) {
	lines := NewLines()

	line := NewSceneLine(0, 1, 10)
	lines.add(line)

	line = NewSceneLine(0, 1, 0)
	lines.add(line)

	line = NewSceneLine(0, 1, 2)
	lines.add(line)

	assert.Equal(t, 0, lines.findIndex(0))
	assert.Equal(t, 2, lines.findIndex(10))
	assert.Equal(t, -1, lines.findIndex(1))
}

func Test_Remove(t *testing.T) {
	lines := NewLines()

	line := NewSceneLine(0, 1, 10)
	lines.add(line)

	line = NewSceneLine(0, 1, 0)
	lines.add(line)

	line = NewSceneLine(0, 1, 2)
	lines.add(line)

	assert.Equal(t, false, lines.findIndex(2) == -1)
	lines.remove(2)
	assert.Equal(t, true, lines.findIndex(2) == -1)

}
