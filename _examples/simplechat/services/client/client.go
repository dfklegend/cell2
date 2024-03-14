package client

import (
	"log"
	"math/rand"
	"time"

	"github.com/asynkron/protoactor-go/actor"

	as "github.com/dfklegend/cell2/actorex/service"
	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	"github.com/dfklegend/cell2/node/app"
	actormodule "github.com/dfklegend/cell2/node/modules/actor"
	"github.com/dfklegend/cell2/node/service"
	mymsg "simplechat/messages"
)

type Service struct {
	*service.NodeService

	connectorId  string
	pidConnector *actor.PID

	logined bool

	pidChatService *actor.PID
	id             int32
	roomID         int32
	Token          int32
	name           string
}

func NewService() *Service {
	s := &Service{
		NodeService: service.NewService(),
		name:        "tom",
	}

	s.Service.InitReqReceiver(s)
	return s
}

func (s *Service) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *service.StartServiceCmd:
		s.NodeService.Receive(ctx)
		s.Start(msg)
		return
	case *mymsg.ClientLogin:
		s.login()
		return
	case *mymsg.ClientSay:
		s.say(msg.Str)
		return
	case *mymsg.ClientNickname:
		s.nickname(msg.Name)
		return
	}
	s.NodeService.Receive(ctx)
}

func (s *Service) ReceiveRequest(ctx actor.Context, request *messages.ServiceRequest, rawMsg interface{}) {
	switch msg := rawMsg.(type) {
	case *mymsg.Chat:
		if s.logined {
			log.Printf("- %v:%v\n", msg.Name, msg.Str)
		}
		s.Response(request, as.CodeSucc, "", nil)
	}
}

func (s *Service) Start(msg *service.StartServiceCmd) {
	log.Printf("client Start: %+v", msg)

	s.GetRunService().GetTimerMgr().AddTimer(4*time.Second, func(args ...interface{}) {
		s.update()
	})
}

func (s *Service) update() {
	if !s.logined {
		return
	}

	if s.pidChatService == nil {
		log.Println("not logined")
		return
	}
	s.Request(s.pidChatService, &mymsg.Ping{
		ID: s.id,
	}, nil)
}

func (s *Service) findGate() *actor.PID {
	nodes := app.Node.GetCluster().GetServiceList("gate")
	if nodes == nil || len(nodes.Items) == 0 {
		return nil
	}
	index := rand.Intn(len(nodes.Items))
	return nodes.Items[index].PID
}

func (s *Service) login() {
	pid := s.findGate()
	if pid == nil {
		log.Println("no gate")
		return
	}

	log.Printf("found gate:%v\n", pid)
	s.Request(pid, &mymsg.ReqConnector{
		X: "123",
	}, func(err error, r interface{}) {
		if err != nil {
			log.Println(err)
			return
		}

		res, _ := r.(*mymsg.AckConnector)
		if res == nil {
			log.Println("bad ret")
			return
		}

		log.Printf("got connector: %v\n", res.ConnectorID)
		s.connectorId = res.ConnectorID
		item := app.Node.GetCluster().GetService(res.ConnectorID)
		if item != nil {
			s.pidConnector = item.PID
		}

		s.loginToConnector()
	})
}

func (s *Service) loginToConnector() {
	s.Request(s.pidConnector, &mymsg.Login{
		Name: "hello",
	}, func(err error, rawRet interface{}) {
		if err != nil {
			log.Printf("login to Connector error: %v\n", err)
			return
		}

		res, _ := rawRet.(*mymsg.AckLogin)
		if res == nil {
			log.Println("bad ret")
			return
		}

		log.Printf("login succ: %v\n", res)
		// continue login to chatservice
		s.loginToChat(res.ID, res.ChatServiceId, res.RoomID, res.Token)
	})
}

func (s *Service) loginToChat(id int32, chatServiceID string, roomID int32, token int32) {
	pid := app.GetServicePID(chatServiceID)
	if pid == nil {
		log.Println("bad chat service id")
		return
	}

	s.Request(pid, &mymsg.ReqJoin{
		Sender: s.Context.Self(),
		ID:     id,
		RoomID: roomID,
		Token:  token,
		Name:   s.name,
	}, func(err error, rawMsg interface{}) {
		if err != nil {
			log.Printf("req join error: %v\n", err)
			return
		}

		s.logined = true
		s.pidChatService = pid
		s.id = id
		s.roomID = roomID
		s.Token = token
	})
}

func (s *Service) say(str string) {
	if s.pidChatService == nil {
		log.Println("not logined")
		return
	}
	s.Request(s.pidChatService, &mymsg.ReqChat{
		ID:  s.id,
		Str: str,
	}, nil)
}

func (s *Service) nickname(str string) {
	s.name = str
	if s.pidChatService == nil {
		log.Println("not logined")
		return
	}
	s.Request(s.pidChatService, &mymsg.Nickname{
		ID:   s.id,
		Name: str,
	}, nil)
}

type Creator struct {
}

func NewCreator() service.IServiceCreator {
	return &Creator{}
}

func (g *Creator) Create(name string) {
	system := actormodule.GetSystem()
	rootContext := system.Root

	props, _ := as.NewServicePropsWithNewScheDisp(func() actor.Actor { return NewService() }, name)
	pid, _ := rootContext.SpawnNamed(props, name)

	service.StartNodeService(rootContext, pid, name, app.Node.GetServiceCfg(name))
}
