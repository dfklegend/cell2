package services

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/actorex/service"
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
	switch msg := ctx.Message().(type) {
	case *mymsgs.Add:
		s.reqAdd(ctx, msg)
	}
	s.Service.Receive(ctx)
}

func (s *ClientService) reqAdd(ctx actor.Context, add *mymsgs.Add) {
	addService1 := actor.NewPID("127.0.0.1:1000", "addService1")

	s.Request(addService1, add, func(err error, r interface{}) {
		if err != nil {
			return
		}
		res, _ := r.(*mymsgs.AddResult)
		log.Printf("result: %v\n", res.Result)
	})
}
