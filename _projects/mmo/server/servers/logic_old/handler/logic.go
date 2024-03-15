package logic

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	builder "github.com/dfklegend/cell2/node/servicebuilder"

	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	"github.com/dfklegend/cell2/node/service"

	"mmo/common/cmd"
	logic2 "mmo/servers/logic_old"
	"mmo/servers/logic_old/cmds"
	"mmo/servers/logic_old/cmds/context"
)

/*
	逻辑服务，玩家基础行为
*/

type Service struct {
	*service.NodeService

	Mgr    *logic2.PlayerMgr
	Retire *logic2.RetireMgr
	State  *logic2.ServiceState

	cmdMgr     *cmd.Mgr
	cmdContext *context.CmdContext
}

func NewService() *Service {
	s := &Service{
		NodeService: service.NewService(),
	}

	s.Mgr = logic2.NewPlayerMgr(s.NodeService)
	s.State = logic2.NewServiceState()
	s.Retire = logic2.NewRetireMgr()

	s.cmdMgr = cmd.NewCmdMgr()
	cmds.RegisterCmds(s.cmdMgr)
	s.cmdContext = &context.CmdContext{
		NS: s.NodeService,
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
}

func (s *Service) Start(msg *service.StartServiceCmd) {
	log.Printf("logic_old Start: %+v", msg)

	s.State.Init(s.Retire)
	s.Retire.Start(s.GetNodeService(), s.Mgr)
	s.SetCtrlCmdListener(s.Retire)
	s.Mgr.Start()
}

func NewCreator() service.IServiceCreator {
	return service.NewFuncCreator(func(name string) {
		builder.StartBackService(name, func() actor.Actor { return NewService() })
	})
}
