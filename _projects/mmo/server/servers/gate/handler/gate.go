package gate

import (
	"log"
	"time"

	"github.com/dfklegend/cell2/node/client/impls"
	"github.com/dfklegend/cell2/node/service"

	"github.com/asynkron/protoactor-go/actor"
	builder "github.com/dfklegend/cell2/node/servicebuilder"

	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	"github.com/dfklegend/cell2/node/app"
)

// gate
// 前端服务
// 处理玩家连接，上下线流程

type Service struct {
	*service.NodeService

	sessions *impls.ClientSessions
}

func NewService() *Service {
	s := &Service{
		NodeService: service.NewService(),
	}

	s.Service.InitReqReceiver(s)
	s.NodeService.SetOwner(s)
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
	log.Printf("client Start: %+v", msg)
	s.initSessions()
}

func (s *Service) randGate() *app.ServiceItem {
	return app.RandGetServiceItem("gate")
}

func (s *Service) allocGate() string {
	item := s.randGate()
	if item == nil {
		return ""
	}
	return item.Name
}

func (s *Service) initSessions() {
	ns := s.NodeService
	s.sessions = ns.GetComponent("sessions").(*impls.SessionsComponent).GetSessions()
	s.sessions.SetOnCloseHandler(onSessionClose)
	s.sessions.SetKickHandler(newKickHandler())

	s.GetRunService().GetTimerMgr().AddTimer(time.Second, s.onUpdate)
}

func (s *Service) onUpdate(args ...any) {
	// TODO: 定期向center发送session激活消息
}

func NewCreator() service.IServiceCreator {
	return service.NewFuncCreator(func(name string) {
		builder.StartFrontService(name, func() actor.Actor { return NewService() })
	})
}
