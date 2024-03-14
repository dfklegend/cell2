package light

import (
	"log"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var gotEvent = false
var eventTimes = 0

func Test_Normal(t *testing.T) {
	ec := NewEventCenter()

	id := ec.Subscribe("event1", func(args ...any) {
		log.Println("got event1")
		gotEvent = true
	})

	gotEvent = false
	ec.Publish("event1")

	assert.Equal(t, true, gotEvent)

	ec.UnsubscribeById("event1", id)
	gotEvent = false
	ec.Publish("event1")
	assert.Equal(t, false, gotEvent)
}

func EventCB(args ...any) {
	gotEvent = true
	log.Println("got event from EventCB")
}

func EventCB1(args ...any) {
	gotEvent = true
	log.Println("got event from EventCB1")
}

type TestClass struct {
	I  int
	CB func()
}

func (t *TestClass) Event(args ...any) {
	eventTimes++
	log.Printf("got event from TestClass.Event{%v %v}\n", reflect.ValueOf(t).Pointer(), t.I)
	if t.CB != nil {
		t.CB()
	}
}

func (t TestClass) Event1(args ...any) {
	eventTimes++
	log.Printf("got event from TestClass.Event1{%v %v}\n", reflect.ValueOf(&t).Pointer(), t.I)
	if t.CB != nil {
		t.CB()
	}
}

func Test_CB(t *testing.T) {
	ec := NewEventCenter()

	ec.Subscribe("event1", EventCB)
	gotEvent = false
	ec.Publish("event1")
	assert.Equal(t, true, gotEvent)

	ec.Unsubscribe("event1", EventCB)
	gotEvent = false
	ec.Publish("event1")
	assert.Equal(t, false, gotEvent)

	ec.Subscribe("event1", EventCB1)
	ec.Publish("event1")

	ec.Unsubscribe("event1", EventCB)
	gotEvent = false
	ec.Publish("event1")
	assert.Equal(t, true, gotEvent)

	ec.Subscribe("event1", EventCB)
	ec.Publish("event1")

	ec.Unsubscribe("event1", EventCB)
	ec.Unsubscribe("event1", EventCB1)

	gotEvent = false
	ec.Publish("event1")
	assert.Equal(t, false, gotEvent)
}

func Test_BaseClass(t *testing.T) {
	ec := NewEventCenter()
	a := &TestClass{}
	a.I = 100

	b := &TestClass{}
	b.I = 200
	ec.Subscribe("event1", a.Event)
	ec.Subscribe("event1", b.Event)

	eventTimes = 0
	ec.Publish("event1")
	// 只注册了一个
	assert.Equal(t, 1, eventTimes)
}

func Test_Class(t *testing.T) {
	ec := NewEventCenter()
	ai := 0
	bi := 0
	a := &TestClass{
		I: 100,
		CB: func() {
			ai++
		}}
	b := &TestClass{
		I: 200,
		CB: func() {
			bi++
		}}

	clearData := func() {
		eventTimes = 0
		ai = 0
		bi = 0
	}

	ec.SubscribeWithReceiver("event1", a, a.Event)
	ec.SubscribeWithReceiver("event1", b, b.Event)

	clearData()
	ec.Publish("event1")
	assert.Equal(t, 2, eventTimes)

	ec.UnsubscribeWithReceiver("event1", b, b.Event)

	clearData()
	ec.Publish("event1")
	assert.Equal(t, 1, eventTimes)
	assert.Equal(t, 1, ai)
	assert.Equal(t, 0, bi)

	// 再次注册
	ec.SubscribeWithReceiver("event1", b, b.Event)
	// 函数是同一个，也能取消掉
	ec.UnsubscribeWithReceiver("event1", b, a.Event)

	clearData()
	ec.Publish("event1")
	assert.Equal(t, 1, eventTimes)
	assert.Equal(t, 1, ai)
	assert.Equal(t, 0, bi)

	ec.SubscribeWithReceiver("event1", b, b.Event)
	ec.UnsubscribeWithReceiver("event1", b, b.Event1)

	clearData()
	ec.Publish("event1")
	assert.Equal(t, 2, eventTimes)
	assert.Equal(t, 1, ai)
	assert.Equal(t, 1, bi)
}

//
func Test_ClassValueReceiver(t *testing.T) {
	ec := NewEventCenter()
	a := &TestClass{I: 100}
	b := &TestClass{I: 200}

	ec.SubscribeWithReceiver("event1", a, a.Event1)
	ec.SubscribeWithReceiver("event1", b, b.Event1)

	eventTimes = 0
	ec.Publish("event1")
	assert.Equal(t, 2, eventTimes)

	ec.UnsubscribeWithReceiver("event1", b, a.Event1)
	eventTimes = 0
	ec.Publish("event1")
	assert.Equal(t, 1, eventTimes)
}
