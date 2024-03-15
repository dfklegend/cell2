package logicservice

/*
logic主体功能
本来命名成service，但是goland一堆提示错误(编译没问题)
*/

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	builder "github.com/dfklegend/cell2/node/servicebuilder"

	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	ns "github.com/dfklegend/cell2/node/service"

	"mmo/common/cmd"
	"mmo/servers/logic_old/cmds"
	"mmo/servers/logic_old/cmds/context"
)

type Service struct {
	*ns.NodeService

	cmdMgr     *cmd.Mgr
	cmdContext *context.CmdContext

	players *PlayerMgr
}

func NewService() *Service {
	s := &Service{
		NodeService: ns.NewService(),
	}

	s.cmdMgr = cmd.NewCmdMgr()
	s.players = NewPlayerMgr(s.NodeService)
	cmds.RegisterCmds(s.cmdMgr)
	s.cmdContext = &context.CmdContext{
		NS: s.NodeService,
	}

	s.Service.InitReqReceiver(s)
	return s
}

func (s *Service) GetNodeService() *ns.NodeService {
	return s.NodeService
}

func (s *Service) GetPlayers() *PlayerMgr {
	return s.players
}

func (s *Service) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *ns.StartServiceCmd:
		s.NodeService.Receive(ctx)
		s.Start(msg)
		return
	}
	s.NodeService.Receive(ctx)
}

func (s *Service) ReceiveRequest(ctx actor.Context, request *messages.ServiceRequest, rawMsg interface{}) {
}

func (s *Service) Start(msg *ns.StartServiceCmd) {
	log.Printf("logic Start: %+v", msg)
	s.players.Start()
}

func NewCreator() ns.IServiceCreator {
	return ns.NewFuncCreator(func(name string) {
		builder.StartBackService(name, func() actor.Actor { return NewService() })
	})
}
