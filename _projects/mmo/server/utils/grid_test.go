package utils

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRound(t *testing.T) {
	i := int(math.Round(0.4))
	assert.Equal(t, 0, i)

	i = int(math.Round(-0.4))
	assert.Equal(t, 0, i)

	i = int(math.Round(-0.6))
	assert.Equal(t, -1, i)

	i = int(math.Round(-1.4))
	assert.Equal(t, -1, i)

	i = int(math.Round(0.9))
	assert.Equal(t, 1, i)

	i = int(math.Round(1.3))
	assert.Equal(t, 1, i)
}

func TestToBlock(t *testing.T) {
	width := 100

	x := GridToX(0, width)
	blockX := ToGridX(x, width)
	assert.Equal(t, true, 0 == blockX)

	blockX = ToGridX(0.5, width)
	assert.Equal(t, 51, blockX)

	blockX = ToGridX(-50.4, width)
	assert.Equal(t, 0, blockX)

	blockX = ToGridX(-49.6, width)
	assert.Equal(t, 0, blockX)

	blockX = ToGridX(-49.4, width)
	assert.Equal(t, 1, blockX)

	blockX = ToGridX(-100, width)
	assert.Equal(t, -50, blockX)

	blockX = ToGridX(-0.4, width)
	assert.Equal(t, 50, blockX)
}
