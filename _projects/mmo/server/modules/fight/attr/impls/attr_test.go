package impls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAttr(t *testing.T) {
	attr := NewAttr()
	attr.SetBase(100)
	attr.OffsetPercent(-0.1)
	assert.Equal(t, 90, attr.GetIntValue())
}

func TestPercentClamper(t *testing.T) {
	attr := NewAttr()
	attr.SetPercentClamper(NewMinPercentClapmper(0.1))

	attr.SetBase(100)
	attr.OffsetPercent(-1.9)
	assert.Equal(t, 10, attr.GetIntValue())
}

func TestSlice(t *testing.T) {
	attrs := NewAttrSlice(100)

	attr := attrs.GetAttr(0)
	attr.SetBase(100)
	attr.OffsetPercent(-0.5)
	assert.Equal(t, 50, attrs.GetAttr(0).GetIntValue())
}
