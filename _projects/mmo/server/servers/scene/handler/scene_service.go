package sceneservice

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	builder "github.com/dfklegend/cell2/node/servicebuilder"
	"github.com/dfklegend/cell2/utils/event/light"

	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	"github.com/dfklegend/cell2/node/service"

	"mmo/common/cmd"
	"mmo/servers/scene"
	"mmo/servers/scene/sceneplayer"
	"mmo/servers/scene/sceneplayer/cmds"
)

/*
	负责玩家战斗
*/

type Service struct {
	*service.NodeService

	cmdMgr  *cmd.Mgr
	Mgr     *sceneservice.SceneMgr
	players *sceneplayer.PlayerMgr
	events  *light.EventCenter
}

func NewService() *Service {
	s := &Service{
		NodeService: service.NewService(),
		events:      light.NewEventCenter(),
	}

	s.cmdMgr = cmd.NewCmdMgr()
	cmds.RegisterCmds(s.cmdMgr)
	s.Mgr = sceneservice.NewSceneMgr()
	s.players = sceneplayer.NewPlayerMgr(s.NodeService)
	s.Service.InitReqReceiver(s)
	s.NodeService.SetOwner(s)
	return s
}

func (s *Service) GetNodeService() *service.NodeService {
	return s.NodeService
}

func (s *Service) GetPlayers() *sceneplayer.PlayerMgr {
	return s.players
}

func (s *Service) GetEvents() *light.EventCenter {
	return s.events
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
	log.Printf("scene Start: %+v", msg)

	s.initLua()
	s.Mgr.Start(s.GetNodeService())
	s.players.Start()
}

func (s *Service) initLua() {
	luaService := initLua()
	s.NodeService.AddComponent("lua", NewLuaComponent(luaService))
	s.Mgr.SetLua(luaService)
}

func NewCreator() service.IServiceCreator {
	return service.NewFuncCreator(func(name string) {
		builder.StartBackService(name, func() actor.Actor { return NewService() })
	})
}
