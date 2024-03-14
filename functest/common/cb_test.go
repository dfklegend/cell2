package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 对象的函数指携带了receiver

type obj struct {
	name string
}

func (o *obj) do() string {
	return o.name
}

func call(doFunc func() string) string {
	return doFunc()
}

func Test_CB(t *testing.T) {
	o1 := &obj{name: "o1"}
	o2 := &obj{name: "o2"}

	assert.Equal(t, "o1", call(o1.do))
	assert.Equal(t, "o2", call(o2.do))
}
