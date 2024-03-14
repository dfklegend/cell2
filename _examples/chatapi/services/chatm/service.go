package chatmgr

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"

	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	"github.com/dfklegend/cell2/node/app"
	actormodule "github.com/dfklegend/cell2/node/modules/actor"
	"github.com/dfklegend/cell2/node/service"
)

type Service struct {
	*service.NodeService
	mgr *Mgr
}

func NewService() *Service {
	s := &Service{
		NodeService: service.NewService(),
		mgr:         NewMgr(),
	}

	s.Service.InitReqReceiver(s)
	return s
}

func (s *Service) GetNodeService() *service.NodeService {
	return s.NodeService
}

func (s *Service) OnCreate() {
	log.Println("chat mgr service onCreate")
	s.mgr.Init(s.GetRunService())
}

func (s *Service) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *service.StartServiceCmd:
		s.NodeService.Receive(ctx)
		s.Start(msg)
		return
		//case *mymsg.CSRefreshInfo:
		//	s.mgr.OnServiceInfo(msg.ChatServiceId, msg.RoomNum, msg.PlayerNum)
		//	return
		//case *mymsg.CSRoomStat:
		//	s.mgr.updateRoomStat(msg.RoomID, msg.PlayerNum)
		//	return
	}
	s.NodeService.Receive(ctx)
}

func (s *Service) ReceiveRequest(ctx actor.Context, request *messages.ServiceRequest, rawMsg interface{}) {
	//switch msg := rawMsg.(type) {
	//case *mymsg.CMReqLogin:
	//	log.Printf("got ReqConnector \n")
	//	s.mgr.Login(s, request, msg)
	//}
}

func (s *Service) Start(msg *service.StartServiceCmd) {
	log.Printf("chatmgr service Start: %+v", msg)
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
		name, "chatm.remote")

	pid, _ := rootContext.SpawnNamed(props, name)

	service.StartNodeService(rootContext, pid, name, app.Node.GetServiceCfg(name))
}
