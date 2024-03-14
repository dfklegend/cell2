package services

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/actorex/service"
	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	"github.com/dfklegend/cell2/utils/common"
	"github.com/dfklegend/cell2/utils/waterfall"

	mymsgs "test-service/messages"
)

type AddService2 struct {
	*service.Service
}

func NewAddService2() *AddService2 {
	s := &AddService2{
		Service: service.NewService(),
	}

	s.Service.InitReqReceiver(s)
	return s
}

func (s *AddService2) Receive(ctx actor.Context) {
	s.Service.Receive(ctx)
}

func (s *AddService2) ReceiveRequest(ctx actor.Context, request *messages.ServiceRequest, rawMsg interface{}) {
	switch msg := rawMsg.(type) {
	case *mymsgs.Add:
		log.Printf("request routine: %v\n", common.GetRoutineID())
		log.Printf("go add: %v\n", msg.I)
		s.add(ctx, request, msg)
	}
}

func (s *AddService2) add(ctx actor.Context, request *messages.ServiceRequest, msg *mymsgs.Add) {

	sche := s.GetRunService().GetScheduler()
	waterfall.Sche(sche, []waterfall.Task{
		func(callback waterfall.Callback, args ...interface{}) {
			go func() {
				time.Sleep(1 * time.Second)
				callback(false)
			}()
		},
		func(callback waterfall.Callback, args ...interface{}) {
			res := &mymsgs.AddResult{
				Result: msg.I + 2,
			}
			s.Response(request, 0, "", res)
		},
	},
		func(error bool, args ...interface{}) {
			//
		})
}
