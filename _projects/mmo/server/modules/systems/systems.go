package systems

import (
	"github.com/dfklegend/cell2/node/service"
	l "github.com/dfklegend/cell2/utils/logger"
)

type Systems struct {
	player ILogicPlayer
	ns     *service.NodeService

	systems   map[string]ISystem
	visitList []ISystem
}

func NewSystems(player ILogicPlayer) *Systems {
	s := &Systems{
		player:    player,
		ns:        player.GetNodeService(),
		visitList: make([]ISystem, 0),
		systems:   make(map[string]ISystem),
	}
	return s
}

func (s *Systems) AddSystem(name string, system ISystem) {

	s.systems[name] = system
	s.visitList = append(s.visitList, system)

	system.Init(s.player)
	system.OnCreate()
}

func (s *Systems) Destroy() {
	s.Visit(func(sys ISystem) {
		sys.OnDestroy()
	})
	s.visitList = make([]ISystem, 0)
	s.systems = make(map[string]ISystem)
}

func (s *Systems) GetSystem(name string) ISystem {
	return s.systems[name]
}

// Visit 按顺序访问系统
func (s *Systems) Visit(visitor func(system ISystem)) {
	for _, v := range s.visitList {
		visitor(v)
	}
}

func (s *Systems) Request(system, cmd string, args []byte, cb func(ret []byte, errCode int32)) {
	sys := s.GetSystem(system)
	if sys == nil {
		l.L.Errorf("can not find %v.%v", system, cmd)
		return
	}
	if !sys.HasCmd(cmd) {
		l.L.Errorf("system :%v has not cmd: %v", system, cmd)
		return
	}

	sys.Request(cmd, args, cb)
}
