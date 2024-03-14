package chat

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"

	as "github.com/dfklegend/cell2/actorex/service"
	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	"github.com/dfklegend/cell2/node/app"
	actormodule "github.com/dfklegend/cell2/node/modules/actor"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/timer"
	"simplechat/common"
	mymsg "simplechat/messages"
)

type Service struct {
	*service.NodeService
	mgr          *RoomMgr
	refreshTimer timer.IdType
}

func NewService() *Service {
	s := &Service{
		NodeService: service.NewService(),
		mgr:         NewMgr(),
	}

	s.mgr.Init(s)
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
	switch msg := rawMsg.(type) {
	case *mymsg.CSCreateRoom:
		log.Printf("got ReqConnector \n")
		s.onCreateRoom(request, msg)
	case *mymsg.ReqJoin:
		s.reqJoin(request, msg)
	case *mymsg.ReqChat:
		s.reqChat(request, msg)
	case *mymsg.Ping:
		s.ping(request, msg)
	case *mymsg.Nickname:
		s.nickname(request, msg)
	}
}

func (s *Service) Start(msg *service.StartServiceCmd) {
	log.Printf("client Start: %+v", msg)

	s.refreshTimer = s.GetRunService().GetTimerMgr().AddTimer(5*time.Second, func(args ...interface{}) {
		s.refreshInfo()
	})
	s.mgr.Start()
}

func (s *Service) refreshInfo() {
	pid := common.GetFirstService("chatm")
	if pid == nil {
		return
	}
	s.Send(pid, &mymsg.CSRefreshInfo{
		ChatServiceId: s.Name,
		RoomNum:       int32(s.mgr.GetRoomNum()),
		PlayerNum:     int32(s.mgr.GetPlayerNum()),
	})
}

func (s *Service) reportRoomStat(roomID int32, playerNum int32) {
	pid := common.GetFirstService("chatm")
	if pid == nil {
		return
	}
	s.Send(pid, &mymsg.CSRoomStat{
		RoomID:    roomID,
		PlayerNum: playerNum,
	})
}

func (s *Service) onCreateRoom(request *messages.ServiceRequest, create *mymsg.CSCreateRoom) {
	if !s.mgr.CreateRoom(create.RoomID, create.Token) {
		s.Response(request, as.CodeErrString, "create room failed!", nil)
		return
	}

	s.Response(request, as.CodeSucc, "", nil)
}

func (s *Service) reqChat(request *messages.ServiceRequest, msg *mymsg.ReqChat) {
	s.mgr.Chat(msg.ID, msg.Str)
	s.Response(request, as.CodeSucc, "", nil)
}

func (s *Service) reqJoin(request *messages.ServiceRequest, msg *mymsg.ReqJoin) {
	if !s.mgr.Join(msg.ID, msg.Sender, msg.Name, msg.RoomID, msg.Token) {
		s.Response(request, as.CodeErrString, "join failed", nil)
		return
	}

	s.Response(request, as.CodeSucc, "", nil)
}

func (s *Service) ping(request *messages.ServiceRequest, msg *mymsg.Ping) {
	s.mgr.Ping(msg.ID)
	s.Response(request, as.CodeSucc, "", nil)
}

func (s *Service) nickname(request *messages.ServiceRequest, msg *mymsg.Nickname) {
	s.mgr.ChangeName(msg.ID, msg.Name)
	s.Response(request, as.CodeSucc, "", nil)
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
