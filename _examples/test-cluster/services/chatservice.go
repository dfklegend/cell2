package services

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

	mymsg "test-cluster/messages"
)

type ChatService struct {
	*service.NodeService
}

func NewChatService() *ChatService {
	s := &ChatService{
		NodeService: service.NewService(),
	}

	s.Service.InitReqReceiver(s)
	return s
}

func (s *ChatService) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *service.StartServiceCmd:
		s.NodeService.Receive(ctx)
		s.Start(msg)
		return
	}
	s.NodeService.Receive(ctx)
}

func (s *ChatService) ReceiveRequest(ctx actor.Context, request *messages.ServiceRequest, rawMsg interface{}) {

}

func (s *ChatService) Start(msg *service.StartServiceCmd) {
	log.Printf("chat Start: %+v", msg)
	s.GetRunService().GetTimerMgr().AddTimer(5*time.Second, func(args ...interface{}) {
		// 随机向gate-1,gate-2发hello
		s.randHello()
	})
}

func (s *ChatService) findGate() *actor.PID {
	// return actor.NewPID("nonhost", "gate-1")
	nodes := app.Node.GetCluster().GetServiceList("gate")
	if nodes == nil || len(nodes.Items) == 0 {
		return nil
	}
	index := rand.Intn(len(nodes.Items))
	return nodes.Items[index].PID
}

func (s *ChatService) randHello() {
	// found a gate service
	pid := s.findGate()
	if pid == nil {
		log.Printf("find no gate\n")
		return
	}

	log.Printf("find gate: %v\n", pid)

	s.Request(pid, &mymsg.ChatHello{
		From: s.Name,
	}, func(err error, r interface{}) {
		res, _ := r.(*mymsg.ChatHelloRet)
		if res == nil {
			return
		}
		log.Printf("got res from :%v\n", res.From)
	})
}

type ChatCreator struct {
}

func NewChatCreator() service.IServiceCreator {
	return &ChatCreator{}
}

func (g *ChatCreator) Create(name string) {
	system := actormodule.GetSystem()
	rootContext := system.Root

	props, _ := as.NewServicePropsWithNewScheDisp(func() actor.Actor { return NewChatService() }, name)
	pid, _ := rootContext.SpawnNamed(props, name)

	service.StartNodeService(rootContext, pid, name, app.Node.GetServiceCfg(name))
}
