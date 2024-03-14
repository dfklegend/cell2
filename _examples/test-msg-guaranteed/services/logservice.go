package services

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/actorex/service"
	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	"github.com/dfklegend/cell2/utils/common"
	l "github.com/dfklegend/cell2/utils/logger"
	mymsgs "test-service/messages"
)

type LogService struct {
	*service.Service
}

func NewAddService() *LogService {
	s := &LogService{
		Service: service.NewService(),
	}

	s.Service.InitReqReceiver(s)
	return s
}

func (s *LogService) Receive(ctx actor.Context) {
	s.Service.Receive(ctx)
}

func (s *LogService) ReceiveRequest(ctx actor.Context, request *messages.ServiceRequest, rawMsg interface{}) {
	switch msg := rawMsg.(type) {
	case *mymsgs.Add:
		log.Printf("request routine: %v\n", common.GetRoutineID())
		s.log(ctx, request, msg)
	}
}

func (s *LogService) log(ctx actor.Context, request *messages.ServiceRequest, inMsg *mymsgs.Add) {
	// just log it
	l.L.Infof("  got %v", inMsg.I)
	s.Response(request, 0, "", nil)
}
