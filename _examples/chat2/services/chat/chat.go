package chat

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"

	"chat2/logics"
	as "github.com/dfklegend/cell2/actorex/service"
	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/client/impls"
	actormodule "github.com/dfklegend/cell2/node/modules/actor"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/timer"
)

type Service struct {
	*service.NodeService
	//mgr          *RoomMgr
	refreshTimer timer.IdType

	//
	cs *logics.ChatMgr
}

func NewService() *Service {
	s := &Service{
		NodeService: service.NewService(),
		//mgr:         NewMgr(),
	}

	//s.mgr.Init(s)
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
	log.Printf("chat Start: %+v", msg)

	ns := s.GetNodeService()
	cs := ns.GetComponent("channel").(*impls.ChannelComponent).GetCS()
	s.cs = logics.NewChatService(ns, cs)

	s.refreshTimer = s.GetRunService().GetTimerMgr().AddTimer(5*time.Second, func(args ...interface{}) {
		//
	})
	//s.mgr.Start()
}

type Creator struct {
}

func NewCreator() service.IServiceCreator {
	return &Creator{}
}

func (g *Creator) Create(name string) {
	system := actormodule.GetSystem()
	rootContext := system.Root

	props, ext := service.NewServiceWithDispatcher(func() actor.Actor { return NewService() },
		name, "chat.remote")
	ext.WithPostFunc(func(s as.IService) {
		ns, _ := s.(service.INodeServiceOwner)
		impls.ServiceCreateCommonComponents(ns.GetNodeService(), app.Node.GetServiceCfg(name))
	})

	pid, _ := rootContext.SpawnNamed(props, name)
	service.StartNodeService(rootContext, pid, name, app.Node.GetServiceCfg(name))
}
