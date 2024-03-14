package connector

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"

	as "github.com/dfklegend/cell2/actorex/service"
	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	"github.com/dfklegend/cell2/node/app"
	actormodule "github.com/dfklegend/cell2/node/modules/actor"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/waterfall"
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
	switch msg := rawMsg.(type) {
	case *mymsg.Login:
		log.Printf("got Login \n")
		s.reqLogin(request, msg)
	}
}

func (s *Service) reqLogin(request *messages.ServiceRequest, login *mymsg.Login) {
	// . 请求loginservice分配一个id
	// . 请求进入ChatMgr
	sche := s.GetRunService().GetScheduler()

	var loginID int32
	waterfall.Sche(sche, []waterfall.Task{
		func(callback waterfall.Callback, args ...interface{}) {
			s.reqLSLogin(login.Name, func(error bool, ackLogin *mymsg.LSAckLogin) {
				if error {
					callback(true)
					return
				}

				loginID = ackLogin.ID
				callback(false)
			})
		},
		func(callback waterfall.Callback, args ...interface{}) {
			s.reqChatMgrLogin(login.Name, loginID, func(error bool, ackLogin *mymsg.CMAckLogin) {
				if error {
					callback(true)
					return
				}
				callback(false, ackLogin)
			})
		},
		func(callback waterfall.Callback, args ...interface{}) {
			ackLogin, _ := args[0].(*mymsg.CMAckLogin)

			s.Response(request, as.CodeSucc, "", &mymsg.AckLogin{
				ID:            loginID,
				ChatServiceId: ackLogin.ChatServiceId,
				RoomID:        ackLogin.RoomID,
				Token:         ackLogin.Token,
			})
			callback(false)
		},
	}, // end tasks
		func(error bool, args ...interface{}) {
			if error {
				s.Response(request, as.CodeErrString, "失败", nil)
			}
		})
}

func (s *Service) reqLSLogin(name string, callback func(error bool, msg *mymsg.LSAckLogin)) {
	pid := common.GetFirstService("login")
	if pid == nil {
		log.Println("no login service")
		callback(true, nil)
		return
	}

	s.Request(pid, &mymsg.LSLogin{
		Name: name,
	}, func(err error, raw interface{}) {
		if err != nil {
			log.Printf("err: %v\n", err)
			callback(true, nil)
			return
		}

		msg, _ := raw.(*mymsg.LSAckLogin)
		if msg == nil {
			callback(true, nil)
			return
		}

		callback(false, msg)
	})
}

func (s *Service) reqChatMgrLogin(name string, ID int32, callback func(error bool, msg *mymsg.CMAckLogin)) {
	pid := common.GetFirstService("chatm")
	if pid == nil {
		log.Println("no chatm service")
		callback(true, nil)
		return
	}

	s.Request(pid, &mymsg.CMReqLogin{
		Name: name,
		ID:   ID,
	}, func(err error, raw interface{}) {
		if err != nil {
			log.Printf("err: %v\n", err)
			callback(true, nil)
			return
		}

		msg, _ := raw.(*mymsg.CMAckLogin)
		if msg == nil {
			callback(true, nil)
			return
		}

		callback(false, msg)
	})
}

func (s *Service) Start(msg *service.StartServiceCmd) {
	log.Printf("connector Start: %+v", msg)
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
