package space

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	space := NewZoneSpace()
	space.Init(-5, -5, 5, 5, 5)
	assert.Equal(t, 0, space.xToZoneX(-10))
	assert.Equal(t, 1, space.xToZoneX(-0))
	assert.Equal(t, 2, space.xToZoneX(100))
}
