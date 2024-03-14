package apientry

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/formater"
)

type TestArg struct {
	Str string `json:"str"`
	I   int    `json:"i"`
}

type TestRet struct {
	Str string `json:"str"`
	I   int    `json:"i"`
}

type TestEntry struct {
	*api.BaseAPIEntry
}

// miss match
func (t *TestEntry) Func1() {
}

func (t *TestEntry) Func2(ctx *api.DummyContext, arg *TestArg, cb func(err error, ret interface{})) {
	log.Printf("enter func2: %v nil:%v\n", arg, ctx == nil)
	cb(nil, &TestRet{
		Str: "ret",
	})
}

// arg not ptr
func (t *TestEntry) Func3(ctx *api.DummyContext, arg TestArg, cb func()) {
}

func (t *TestEntry) NotifyFunc2(ctx *api.DummyContext, arg *TestArg) {
	log.Printf("enter NotifyFunc2: %v nil:%v\n", arg, ctx == nil)
}

func TestAnalysis(t *testing.T) {
	e := &TestEntry{}
	c := NewContainer(e)
	c.ExtractHandler(formater.GetDefaultFormater())

	assert.Equal(t, false, c.HasMethod("Func1"))
	assert.Equal(t, true, c.HasMethod("Func2"))
	assert.Equal(t, false, c.HasMethod("Func3"))
	assert.Equal(t, true, c.HasMethod("NotifyFunc2"))
}

func TestCall(t *testing.T) {
	e := &TestEntry{}
	c := NewContainer(e)
	c.ExtractHandler(formater.GetDefaultFormater())

	inArg := &TestArg{
		Str: "Hello",
	}

	ret1 := false
	c.CallMethod(nil, "Func2", inArg, func(err error, ret interface{}) {
		log.Printf("ret1: %v\n", ret)
		ret1 = true
	})

	ret2 := false
	c.CallMethod(&api.DummyContext{}, "Func2", inArg, func(err error, ret interface{}) {
		log.Printf("ret2: %v\n", ret)
		ret2 = true
	})

	ret3 := false
	c.CallMethod(&api.DummyContext{}, "NotifyFunc2", inArg, func(err error, ret interface{}) {
		log.Printf("ret3: %v\n", ret)
		ret3 = true
	})

	c.CallMethod(&api.DummyContext{}, "NotifyFunc2", inArg, nil)

	c.CallMethod(&api.DummyContext{}, "NotifyFunc3", inArg, func(err error, ret interface{}) {
	})

	time.Sleep(time.Millisecond)
	assert.Equal(t, true, ret1)
	assert.Equal(t, true, ret2)
	assert.Equal(t, false, ret3)
}
