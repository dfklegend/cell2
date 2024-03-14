package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AutoRound(t *testing.T) {
	s := &SerialIdService{}

	assert.Equal(t, uint32(1), s.AllocId())
	s.nextId = 0xFFFFFFFE
	assert.Equal(t, uint32(0xFFFFFFFF), s.AllocId())

	assert.Equal(t, uint32(0), s.AllocId())
	assert.Equal(t, uint32(1), s.AllocId())
}
