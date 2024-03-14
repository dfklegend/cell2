package gate

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"

	as "github.com/dfklegend/cell2/actorex/service"
	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	"github.com/dfklegend/cell2/node/app"
	actormodule "github.com/dfklegend/cell2/node/modules/actor"
	"github.com/dfklegend/cell2/node/service"
	"simplechat/common"
	mymsg "simplechat/messages"
)

type Service struct {
	*service.NodeService
}

func NewService() *Service {
	s := &Service{
		NodeService: service.NewService(),
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
	case *mymsg.ReqConnector:
		log.Printf("got ReqConnector \n")
		s.reqConnector(request)
	}
}

func (s *Service) Start(msg *service.StartServiceCmd) {
	log.Printf("client Start: %+v", msg)
}

func (s *Service) findConnector() *app.ServiceItem {
	return common.RandGetServiceItem("connector")
}

func (s *Service) allocConnector() string {
	item := s.findConnector()
	if item == nil {
		return ""
	}
	return item.Name
}

func (s *Service) reqConnector(request *messages.ServiceRequest) {
	id := s.allocConnector()

	if id == "" {
		s.Response(request, as.CodeErrString, "没有合适的connector", nil)
		return
	}
	s.Response(request, 0, "", &mymsg.AckConnector{
		ConnectorID: id,
	})
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
