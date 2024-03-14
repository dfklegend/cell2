package service

import (
	"log"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"

	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	"github.com/dfklegend/cell2/utils/common"

	"github.com/stretchr/testify/assert"

	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/utils/serialize/proto"
	"github.com/dfklegend/cell2/utils/serialize/proto/msgs"
)

type EntryProto struct {
	api.APIEntry
}

func (e *EntryProto) Join(d *RemoteContext, msg *msgs.TestHello, cbFunc apientry.HandlerCBFunc) error {
	log.Printf("in EntryProto.join\n")
	log.Printf("join routine: %v\n", common.GetRoutineID())

	apientry.CheckInvokeCBFunc(cbFunc, nil, &msgs.TestHello{
		I: msg.I * 2,
	})
	theApiTestEnv.joined++
	return nil
}

func (e *EntryProto) Notify(d *RemoteContext, msg *msgs.TestHello) error {
	log.Printf("in EntryProto.Notify\n")
	log.Printf("Notify routine: %v\n", common.GetRoutineID())
	theApiTestEnv.notified++
	return nil
}

type TestCall struct {
	pid    *actor.PID
	route  string
	notify bool
}

type apiTestEnv struct {
	result   int32
	err      error
	joined   int32
	notified int32
}

func (t *apiTestEnv) Reset() {
	t.result = 0
	t.err = nil
	t.joined = 0
	t.notified = 0
}

var theApiTestEnv = &apiTestEnv{}

type APIHelloService struct {
	*Service
}

func NewAPIHelloService() *APIHelloService {
	s := &APIHelloService{
		Service: NewService(),
	}

	s.Service.InitReqReceiver(s)
	return s
}

func (s *APIHelloService) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *TestCall:
		if msg.notify {
			s.TestNotify(msg.pid, msg.route)
		} else {
			s.TestCall(msg.pid, msg.route)
		}
		return
	}
	s.Service.Receive(ctx)
}

func (s *APIHelloService) ReceiveRequest(ctx actor.Context, request *messages.ServiceRequest, rawMsg interface{}) {
}

//	模拟逻辑

func (s *APIHelloService) TestCall(pid *actor.PID, route string) {
	hello := &msgs.TestHello{
		I: 100,
	}
	log.Printf("call routine: %v\n", common.GetRoutineID())
	log.Printf("call %v\n", route)

	s.RequestEx(pid, route, hello, func(err error, rawMsg interface{}) {
		if err != nil {
			log.Printf("got error: %v\n", err)
			theApiTestEnv.err = err
			return
		}
		msg, ok := rawMsg.(*msgs.TestHello)
		if !ok {
			log.Printf("got response no ret \n")
			return
		}
		log.Printf("got response in cb %v \n", msg.I)
		log.Printf("cb routine: %v\n", common.GetRoutineID())
		theApiTestEnv.result += msg.I
	})
}

func (s *APIHelloService) TestNotify(pid *actor.PID, route string) {
	hello := &msgs.TestHello{
		I: 100,
	}
	log.Printf("call routine: %v\n", common.GetRoutineID())
	log.Printf("call %v\n", route)

	s.NotifyEx(pid, route, hello)
}

func Test_CallProtoRPC(t *testing.T) {
	serializer := proto.GetDefaultSerializer()
	col := apientry.NewCollection()
	col.Register(&EntryProto{}, apientry.WithGroupName("hello"),
		apientry.WithSerializer(serializer), apientry.WithSerializeRet(false)).
		Build()

	system := actor.NewActorSystem()

	// define root Context
	rootContext := system.Root

	props, ext := NewServicePropsWithNewScheDisp(func() actor.Actor { return NewAPIHelloService() }, "")
	ext.WithDispatcher(NewDispatcher(col))
	//

	theEnv.Reset()
	client1, _ := rootContext.SpawnNamed(props, "client1")
	client2, _ := rootContext.SpawnNamed(props, "client2")

	theEnv.Reset()
	rootContext.Send(client1, &TestCall{
		pid:   client2,
		route: "hello.Join",
	})

	time.Sleep(2 * time.Second)
	assert.Equal(t, int32(200), theApiTestEnv.result)
	system.Shutdown()
}

func TestAddColletion(t *testing.T) {
	col := apientry.NewCollection()
	d := NewDispatcher(col, nil, col)
	assert.Equal(t, 2, len(d.apis))
	d = NewDispatcher()
	d.AddCollection(col, nil, col)
	assert.Equal(t, 2, len(d.apis))
}

func Test_CallMultiCols(t *testing.T) {
	col := apientry.NewCollection()
	col.Register(&EntryProto{}, apientry.WithGroupName("hello")).
		Register(&EntryProto{}, apientry.WithGroupName("hello1")).
		Build()

	system := actor.NewActorSystem()

	// define root Context
	rootContext := system.Root

	props, ext := NewServicePropsWithNewScheDisp(func() actor.Actor { return NewAPIHelloService() }, "")
	ext.WithDispatcher(NewDispatcher(col))
	//

	theApiTestEnv.Reset()
	client1, _ := rootContext.SpawnNamed(props, "client1")

	props, ext = NewServicePropsWithNewScheDisp(func() actor.Actor { return NewAPIHelloService() }, "")
	ext.WithDispatcher(NewDispatcher(col))
	client2, _ := rootContext.SpawnNamed(props, "client2")

	log.Printf(" -- call hello.Join\n")
	rootContext.Send(client1, &TestCall{
		pid:   client2,
		route: "hello.Join",
	})
	time.Sleep(1000 * time.Millisecond)
	assert.Equal(t, int32(200), theApiTestEnv.result)
	assert.Equal(t, int32(1), theApiTestEnv.joined)
	theApiTestEnv.Reset()

	log.Printf(" -- call hello1.Join\n")
	rootContext.Send(client1, &TestCall{
		pid:   client2,
		route: "hello1.Join",
	})
	time.Sleep(1000 * time.Millisecond)
	assert.Equal(t, int32(200), theApiTestEnv.result)
	assert.Equal(t, int32(1), theApiTestEnv.joined)
	theApiTestEnv.Reset()

	log.Printf(" -- call hello1.Join no cb\n")
	rootContext.Send(client1, &TestCall{
		pid:    client2,
		route:  "hello1.Join",
		notify: true,
	})
	time.Sleep(1000 * time.Millisecond)
	assert.Equal(t, int32(1), theApiTestEnv.joined)
	theApiTestEnv.Reset()

	// should failed
	log.Printf(" -- call hello1.Notify with cb\n")
	rootContext.Send(client1, &TestCall{
		pid:   client2,
		route: "hello1.Notify",
	})
	time.Sleep(1000 * time.Millisecond)
	assert.Equal(t, int32(0), theApiTestEnv.notified)
	theApiTestEnv.Reset()

	log.Printf(" -- call hello1.Notify normally\n")
	rootContext.Send(client1, &TestCall{
		pid:    client2,
		notify: true,
		route:  "hello1.Notify",
	})
	time.Sleep(1000 * time.Millisecond)
	assert.Equal(t, int32(0), theApiTestEnv.result)
	assert.Equal(t, int32(1), theApiTestEnv.notified)
	theApiTestEnv.Reset()

	// miss
	log.Printf(" -- call hello2.join, will miss\n")
	rootContext.Send(client1, &TestCall{
		pid:   client2,
		route: "hello2.Join",
	})
	time.Sleep(1000 * time.Millisecond)
	assert.Equal(t, int32(0), theApiTestEnv.joined)
	theApiTestEnv.Reset()
}
