package interfacetest

import (
	"log"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type IA interface {
	A()
}

type IB interface {
	B()
}

type ClassAB struct {
}

func (c *ClassAB) A() {
}

func (c *ClassAB) B() {
}

// 接口之间互转
func Test_InterfaceTransform(t *testing.T) {
	one := &ClassAB{}
	var a IA
	var some interface{}

	some = one
	a1 := some.(IA)
	a = one

	//a2 := one.(IA) // 报错 类型断言必须从interface{}调用
	a2 := interface{}(one).(IA)

	assert.Equal(t, true, a != nil)
	b := a.(IB)
	assert.Equal(t, true, b != nil)
	assert.Equal(t, true, a1 != nil)
	assert.Equal(t, true, a2 != nil)

	log.Printf("%v, %v\n", reflect.ValueOf(a).Pointer(), reflect.ValueOf(b).Pointer())
}
