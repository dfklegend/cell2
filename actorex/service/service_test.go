package service

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/stretchr/testify/assert"

	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	"github.com/dfklegend/cell2/utils/common"
)

type TestOne struct {
	pid *actor.PID
	i   int32
}

type TestEnv struct {
	result int32
}

func (t *TestEnv) Reset() {
	t.result = 0
}

var theEnv = &TestEnv{}

type HelloService struct {
	*Service
}

func NewHelloService() *HelloService {
	s := &HelloService{
		Service: NewService(),
	}

	s.Service.InitReqReceiver(s)
	return s
}

func (s *HelloService) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *TestOne:
		s.TestOne(msg.pid, msg.i)
		return
	}
	s.Service.Receive(ctx)
}

func (s *HelloService) ReceiveRequest(ctx actor.Context, request *messages.ServiceRequest, rawMsg interface{}) {
	switch msg := rawMsg.(type) {
	case *messages.TestHello:
		// 直接扔回去
		res := &messages.TestHello{
			I: msg.I * 2,
		}
		s.Response(request, 0, "", res)
		log.Printf("receive Request routine: %v\n", common.GetRoutineID())
	}
}

//	模拟逻辑

func (s *HelloService) TestOne(pid *actor.PID, i int32) {
	fmt.Printf("got test one\n")
	hello := &messages.TestHello{
		I: i,
	}
	s.Request(pid, hello, func(err error, rawMsg interface{}) {
		if err != nil {
			log.Printf("got error: %v\n", err)
			return
		}
		msg, ok := rawMsg.(*messages.TestHello)
		if !ok {
			return
		}
		log.Printf("got response in cb %v \n", msg.I)
		log.Printf("cb routine: %v\n", common.GetRoutineID())
		theEnv.result += msg.I
	})

	rs := s.GetRunService()
	rs.GetTimerMgr().After(100*time.Millisecond, func(args ...interface{}) {
		log.Printf("timer reach!")
	})
}

// 简单call测试
func Test_Call(t *testing.T) {
	system := actor.NewActorSystem()

	// define root Context
	rootContext := system.Root

	//disp := dispp.NewScheDisp()
	//disp.Start()
	//props := actor.PropsFromProducer(func() actor.Actor { return NewHelloService() },
	//	actor.WithDispatcher(disp))
	props, _ := NewServicePropsWithNewScheDisp(func() actor.Actor { return NewHelloService() }, "")
	//

	theEnv.Reset()
	client1, _ := rootContext.SpawnNamed(props, "client1")
	client2, _ := rootContext.SpawnNamed(props, "client2")

	rootContext.Send(client1, &TestOne{
		pid: client2,
		i:   100,
	})

	time.Sleep(1 * time.Second)
	assert.Equal(t, int32(200), theEnv.result)
}

func Test_Timeout(t *testing.T) {
	testReqTimeoutEnable = true
	system := actor.NewActorSystem()

	// define root Context
	rootContext := system.Root

	//disp := dispp.NewScheDisp()
	//disp.Start()
	//props := actor.PropsFromProducer(func() actor.Actor { return NewHelloService() },
	//	actor.WithDispatcher(disp))
	props, _ := NewServicePropsWithNewScheDisp(func() actor.Actor { return NewHelloService() }, "")
	//

	theEnv.Reset()
	client1, _ := rootContext.SpawnNamed(props, "client1")
	client2, _ := rootContext.SpawnNamed(props, "client2")

	for i := 0; i < 5; i++ {
		rootContext.Send(client1, &TestOne{
			pid: client2,
			i:   100 + int32(i),
		})
	}

	time.Sleep(5 * time.Second)
	testReqTimeoutEnable = false
	assert.Equal(t, int32(0), theEnv.result)
}
