package apientry

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/dfklegend/cell2/utils/serialize/json"
	"github.com/dfklegend/cell2/utils/serialize/proto"

	"github.com/stretchr/testify/assert"

	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/formater"
	"github.com/dfklegend/cell2/utils/serialize/proto/msgs"
)

type InMsg struct {
	Abc string `json:"abc"`
}

type OutMsg struct {
	Def string `json:"def"`
}

type Entry1 struct {
	api.APIEntry
}

func (e *Entry1) Join(d *api.DummyContext, msg *InMsg, cbFunc HandlerCBFunc) {
	log.Printf("in join\n")

	CheckInvokeCBFunc(cbFunc, nil, &OutMsg{"dddd"})
}

type Entry2 struct {
	api.APIEntry
}

func (e *Entry2) Join(d *api.DummyContext, msg *InMsg, cbFunc HandlerCBFunc) {
	log.Printf("in join\n")

	CheckInvokeCBFunc(cbFunc, nil, &OutMsg{"dddd"})
}

func Test_CreateContainer(t *testing.T) {
	entry1 := NewContainer(&Entry1{}, WithGroupName("hello"))
	entry1.ExtractHandler(formater.GetDefaultFormater())

	assert.Equal(t, false, entry1.HasMethod("join"), "error groupName")

	entry1 = NewContainer(&Entry1{}, WithNameFunc(strings.ToLower))
	entry1.ExtractHandler(formater.GetDefaultFormater())

	assert.Equal(t, true, entry1.HasMethod("join"), "error groupName")

}

func Test_CreateContainer2(t *testing.T) {
	// entry2 := NewContainer(&Entry2{}, WithNameFunc(strings.ToLower))
	// entry2.ExtractHandler(&HandlerFormater{})

	// assert.Equal(t, true, entry2.HasMethod("join"), "error groupName")
}

func Test_CreateCollection(t *testing.T) {
	col := NewCollection()
	col.Register(&Entry1{}, WithGroupName("hello")).
		Register(&Entry1{}, WithNameFunc(strings.ToLower)).
		Build()

	assert.Equal(t, true, col.HasMethod("hello.Join"), "error groupName")
	assert.Equal(t, true, col.HasMethod("entry1.join"), "error groupName")
}

func Test_CallRPC(t *testing.T) {
	serializer := json.GetDefaultSerializer()
	col := NewCollection()
	col.Register(&Entry1{}, WithGroupName("hello")).
		Register(&Entry1{}, WithNameFunc(strings.ToLower)).
		Build()

	CallWithSerialize(col, nil,
		"hello.Join",
		[]byte("ddd"),
		func(e error, result interface{}) {
			log.Printf("got result:%+v\n", result)
		}, serializer)

	data, _ := serializer.Marshal(&InMsg{
		Abc: "some",
	})

	CallWithSerialize(col, nil,
		"hello.Join",
		data,
		func(e error, result interface{}) {
			log.Printf("got result:%+v\n", result)
		}, serializer)
	//col.Call(nil, "hello.Join", []byte("ddd"), nil)
}

type EntryProto struct {
	api.APIEntry
}

func (e *EntryProto) Join(d *api.DummyContext, msg *msgs.TestHello, cbFunc HandlerCBFunc) {
	log.Printf("in EntryProto.join\n")

	CheckInvokeCBFunc(cbFunc, nil, &msgs.TestHello{
		I: 99,
		S: "Hello",
	})
}

func Test_CallProtoRPC(t *testing.T) {
	serializer := proto.GetDefaultSerializer()
	col := NewCollection()
	col.Register(&EntryProto{}, WithGroupName("hello"), WithSerializer(serializer)).
		Build()

	inBytes, _ := serializer.Marshal(&msgs.TestHello{
		I: 100,
		S: "in",
	})

	finalresult := &msgs.TestHello{}
	CallWithSerialize(
		col,
		nil,
		"hello.Join",
		inBytes,
		func(e error, result interface{}) {
			log.Printf("got result:%+v\n", result)
			f, _ := result.(*msgs.TestHello)
			finalresult.I = f.I
		}, serializer)

	time.Sleep(time.Second)
	assert.Equal(t, int32(99), finalresult.I)
}
