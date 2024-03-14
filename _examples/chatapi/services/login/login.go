package login

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"

	mymsg "chatapi/messages"
	as "github.com/dfklegend/cell2/actorex/service"
	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	"github.com/dfklegend/cell2/node/app"
	actormodule "github.com/dfklegend/cell2/node/modules/actor"
	"github.com/dfklegend/cell2/node/service"
)

type Service struct {
	*service.NodeService

	nextID int32
}

func NewService() *Service {
	s := &Service{
		NodeService: service.NewService(),
		nextID:      1,
	}

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
	switch rawMsg.(type) {
	case *mymsg.LSLogin:
		log.Printf("got LSLogin \n")

		id := s.nextID
		s.nextID++
		s.Response(request, as.CodeSucc, "", &mymsg.LSAckLogin{
			ID: id,
		})
	}
}

func (s *Service) Start(msg *service.StartServiceCmd) {
	log.Printf("client Start: %+v", msg)
}

type Creator struct {
}

func NewCreator() service.IServiceCreator {
	return &Creator{}
}

func (g *Creator) Create(name string) {
	system := actormodule.GetSystem()
	rootContext := system.Root

	props, _ := service.NewServiceWithDispatcher(func() actor.Actor { return NewService() },
		name, "")
	pid, _ := rootContext.SpawnNamed(props, name)

	service.StartNodeService(rootContext, pid, name, app.Node.GetServiceCfg(name))
}
