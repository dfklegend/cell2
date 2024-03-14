package services

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/actorex/service"
	l "github.com/dfklegend/cell2/utils/logger"
	mymsgs "test-service/messages"
)

type ClientService struct {
	*service.Service
}

func NewClientService() *ClientService {
	s := &ClientService{
		Service: service.NewService(),
	}

	s.Service.InitReqReceiver(s)
	return s
}

func (s *ClientService) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *mymsgs.Add:
		s.startTest(ctx)
	}
	s.Service.Receive(ctx)
}

func (s *ClientService) startTest(ctx actor.Context) {

	i := 0
	s.GetRunService().GetTimerMgr().AddTimer(time.Second, func(args ...interface{}) {
		i++
		s.send(ctx, int32(i))
	})
}

func (s *ClientService) send(ctx actor.Context, i int32) {
	addService1 := actor.NewPID("127.0.0.1:1000", "addService1")

	l.L.Infof("  send: %v", i)
	s.Request(addService1, &mymsgs.Add{
		I: i,
	}, func(err error, r interface{}) {
	})
}
