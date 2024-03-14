package nodectrl

import (
	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/actorex/service"
	"github.com/dfklegend/cell2/nodectrl/define"
)

// AdminService
// 接受admin服务器的命令
type AdminService struct {
	*service.Service
	ctrl *NodeCtrl
}

func newAdminService() *AdminService {
	s := &AdminService{
		Service: service.NewService(),
	}

	s.Service.InitReqReceiver(s)
	return s
}

func CreateAdminService(system *actor.ActorSystem, postFunc func(s service.IService)) *actor.PID {
	rootContext := system.Root

	name := define.NodeAdmin
	props, ext := service.NewServicePropsWithNewScheDisp(func() actor.Actor { return newAdminService() },
		name)
	ext.WithAPIs("node.admin")
	ext.WithPostFunc(postFunc)
	pid, _ := rootContext.SpawnNamed(props, name)
	return pid
}

func (s *AdminService) SetCtrl(ctrl *NodeCtrl) {
	s.ctrl = ctrl
}

func (s *AdminService) ProcessCmd(cmd string) string {
	return s.ctrl.ProcessCmd(cmd)
}

func (s *AdminService) ProcessServiceCmd(name, cmd string) string {
	return s.ctrl.ProcessServiceCmd(name, cmd)
}
