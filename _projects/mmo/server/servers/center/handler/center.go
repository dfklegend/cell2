package handler

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	builder "github.com/dfklegend/cell2/node/servicebuilder"

	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	"github.com/dfklegend/cell2/node/service"

	"mmo/servers/center"
)

/*
	负责维护某个玩家在那个gate上，顶人功能
	目前唯一，可以扩展成多个，基于uid来分配
*/

type Service struct {
	*service.NodeService
	Mgr *center.PlayerMgr
}

func NewService() *Service {
	s := &Service{
		NodeService: service.NewService(),
	}

	s.Mgr = center.NewPlayerMgr(s.NodeService)
	s.Service.InitReqReceiver(s)
	return s
}

func (s *Service) GetNodeService() *service.NodeService {
	return s.NodeService
}

func (s *Service) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *service.StartServiceCmd:
		s.NodeService.Receive(ctx)
		s.Start(msg)
		return
	}
	s.NodeService.Receive(ctx)
}

func (s *Service) ReceiveRequest(ctx actor.Context, request *messages.ServiceRequest, rawMsg interface{}) {
}

func (s *Service) Start(msg *service.StartServiceCmd) {
	log.Printf("handler Start: %+v", msg)

	s.Mgr.Start()
}

func NewCreator() service.IServiceCreator {
	return service.NewFuncCreator(func(name string) {
		builder.StartBackService(name, func() actor.Actor { return NewService() })
	})
}
