package services

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"

	as "github.com/dfklegend/cell2/actorex/service"
	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	"github.com/dfklegend/cell2/node/app"
	actormodule "github.com/dfklegend/cell2/node/modules/actor"
	"github.com/dfklegend/cell2/node/service"
	mymsg "test-app/messages"
)

type GateService struct {
	*service.NodeService
}

func NewGateService() *GateService {
	s := &GateService{
		NodeService: service.NewService(),
	}

	s.Service.InitReqReceiver(s)
	return s
}

func (s *GateService) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *service.StartServiceCmd:
		s.NodeService.Receive(ctx)
		s.Start(msg)
		return
	}
	s.NodeService.Receive(ctx)
}

func (s *GateService) ReceiveRequest(ctx actor.Context, request *messages.ServiceRequest, rawMsg interface{}) {
	switch msg := rawMsg.(type) {
	case *mymsg.ChatHello:
		log.Printf("got chat: %v\n", msg.From)
		s.Response(request, 0, "", &mymsg.ChatHelloRet{
			From: s.Name,
		})
	}
}

func (s *GateService) Start(msg *service.StartServiceCmd) {
	log.Printf("gate Start: %+v", msg)
	args := msg.Args
	if args == nil {
		return
	}
	if len(args) != 1 {
		return
	}
	timePrint := (time.Duration)(args[0].(int))
	s.GetRunService().GetTimerMgr().AddTimer(timePrint*time.Millisecond, func(args ...interface{}) {
		log.Printf("%v tick\n", s.Name)
	})
}

type GateCreator struct {
}

func NewGateCreator() service.IServiceCreator {
	return &GateCreator{}
}

func (g *GateCreator) Create(name string) {
	system := actormodule.GetSystem()
	rootContext := system.Root

	props, _ := as.NewServicePropsWithNewScheDisp(func() actor.Actor { return NewGateService() }, name)
	pid, _ := rootContext.SpawnNamed(props, name)

	service.StartNodeService(rootContext, pid, name, app.Node.GetServiceCfg(name), 2000)
}
