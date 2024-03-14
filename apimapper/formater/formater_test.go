package formater

import (
	"log"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	api "github.com/dfklegend/cell2/apimapper"
)

type TestEntry struct {
}

// 参数不一致
func (t *TestEntry) Func1() {
}

// 合适的
func (t *TestEntry) Func2(ctx *api.DummyContext, entry *TestEntry, cb func()) {
}

// entry不是指针
func (t *TestEntry) Func3(ctx *api.DummyContext, entry TestEntry, cb func()) {
}

func Test_Check(t *testing.T) {
	entry := &TestEntry{}

	ty := reflect.TypeOf(entry)

	for m := 0; m < ty.NumMethod(); m++ {
		method := ty.Method(m)

		log.Printf("%v: %v",
			method.Name,
			defaultFormater.IsValidMethod(method))
	}

	// false
	assert.Equal(t, false, defaultFormater.IsValidMethod(ty.Method(0)))
	assert.Equal(t, true, defaultFormater.IsValidMethod(ty.Method(1)))
	assert.Equal(t, false, defaultFormater.IsValidMethod(ty.Method(2)))
}
