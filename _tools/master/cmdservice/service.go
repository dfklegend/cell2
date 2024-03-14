package cmdservice

import (
	"master/cmdservice/cmd"
	"master/cmdservice/cmds"

	"github.com/dfklegend/cell2/actorex/service"
)

type Service struct {
	*service.Service
	cmds *cmd.Mgr
}

func NewService() *Service {
	s := &Service{
		Service: service.NewService(),
		cmds:    cmd.NewMgr(),
	}

	s.Service.InitReqReceiver(s)
	return s
}

func (s *Service) Start() {
	cmds.InitCmds(s.cmds)
}

func (s *Service) ProcessCmd(cmd string, cb func(string)) {
	s.cmds.Call(s.Service, cmd, cb)
}
