package redis

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dfklegend/cell2/actorex/service"
	l "github.com/dfklegend/cell2/utils/logger"
	goredis "github.com/go-redis/redis/v7"

	"mmo/servers/db/dbop"
)

// Service
// 本地redis logicservice，异步访问
// 认证部分
type Service struct {
	*service.Service

	Address string
	Client  *goredis.Client
	IsReady bool
}

func newRedisService() *Service {
	s := &Service{
		Service: service.NewService(),
	}

	s.Service.InitReqReceiver(s)
	return s
}

func CreateRedisService(ownerName string, system *actor.ActorSystem) *actor.PID {
	rootContext := system.Root

	name := fmt.Sprintf("%v.handler.redis", ownerName)
	props, ext := service.NewServicePropsWithNewScheDisp(func() actor.Actor { return newRedisService() },
		name)
	ext.WithAPIs("db.redis")
	pid, _ := rootContext.SpawnNamed(props, name)
	return pid
}

func (s *Service) InitRedis(address string) {
	l.Log.Infof("init redis logicservice: %v", address)
	s.Address = address
	s.Client = goredis.NewClient(&goredis.Options{
		Addr:     address,
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: 1,  // 只允许一个连接，角色load&save有序
	})

	s.checkReady()
}

func (s *Service) checkReady() {
	_, err := s.Client.Ping().Result()
	if err != nil {
		l.L.Errorf("start redis logicservice error: %v", err)
	} else {
		l.L.Infof("redis logicservice is ready: %v", s.Address)
		if dbop.TryInitDB(s.Client) == nil {
			l.L.Infof("redis logicservice init handler succ", s.Address)
			s.IsReady = true
		}
	}
}

func (s *Service) TryCheckReady() bool {
	if s.IsReady {
		return true
	}
	s.checkReady()
	return s.IsReady
}
