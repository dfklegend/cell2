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

type AddService struct {
	*service.Service
}

func NewAddService() *AddService {
	s := &AddService{
		Service: service.NewService(),
	}

	s.Service.InitReqReceiver(s)
	return s
}

func (s *AddService) Receive(ctx actor.Context) {
	s.Service.Receive(ctx)
}

func (s *AddService) ReceiveRequest(ctx actor.Context, request *messages.ServiceRequest, rawMsg interface{}) {
	switch msg := rawMsg.(type) {
	case *mymsgs.Add:
		log.Printf("request routine: %v\n", common.GetRoutineID())
		s.add(ctx, request, msg)
	}
}

func (s *AddService) add(ctx actor.Context, request *messages.ServiceRequest, inMsg *mymsgs.Add) {

	addService2 := actor.NewPID("127.0.0.1:1001", "addService2")

	sche := s.GetRunService().GetScheduler()
	waterfall.Sche(sche, []waterfall.Task{
		func(callback waterfall.Callback, args ...interface{}) {
			log.Printf("routine: %v\n", common.GetRoutineID())
			go func() {
				time.Sleep(1 * time.Second)
				callback(false)
			}()
		},
		func(callback waterfall.Callback, args ...interface{}) {
			log.Printf("routine: %v\n", common.GetRoutineID())
			msg := &mymsgs.Add{
				I: inMsg.I + 1,
			}
			// 请求其他服务器
			s.Request(addService2, msg, func(err error, r interface{}) {
				log.Printf("routine: %v\n", common.GetRoutineID())
				log.Printf("got result: %v\n", r)
				if err != nil {
					callback(true)
					return
				}
				res, ok := r.(*mymsgs.AddResult)
				if !ok {
					callback(true)
					return
				}

				s.Response(request, 0, "", res)
				callback(false)
			})
		},
		func(callback waterfall.Callback, args ...interface{}) {
			callback(false)
		},
	},
		func(error bool, args ...interface{}) {
			//
		})
}
