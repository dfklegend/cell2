package handler

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	actormodule "github.com/dfklegend/cell2/node/modules/actor"
	"github.com/dfklegend/cell2/node/service"
	builder "github.com/dfklegend/cell2/node/servicebuilder"

	"mmo/common/config"
	mymsg "mmo/messages"
	"mmo/servers/db/redisservice"
)

type Service struct {
	*service.NodeService

	// local logicservice
	redis *actor.PID
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
}

func (s *Service) Start(msg *service.StartServiceCmd) {
	log.Printf("handler Start: %+v", msg)
	s.initRedis()
}

func (s *Service) initRedis() {
	info := config.Cfg.DB

	s.redis = redis.CreateRedisService(s.Name, actormodule.GetSystem())

	s.RequestEx(s.redis, "redis.init", &mymsg.DBRedisInit{
		Address: info.Redis,
	}, nil)
}

func (s *Service) initDB() {

}

func NewCreator() service.IServiceCreator {
	return service.NewFuncCreator(func(name string) {
		builder.StartBackService(name, func() actor.Actor { return NewService() })
	})
}
