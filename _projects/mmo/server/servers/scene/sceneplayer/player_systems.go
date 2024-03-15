package sceneplayer

import (
	l "github.com/dfklegend/cell2/utils/logger"

	"mmo/common/factory"
	mymsg "mmo/messages"
	"mmo/modules/systems"
	sysfactory "mmo/servers/scene/sceneplayer/systems/factory"
	"mmo/servers/scene/sceneplayer/systems/interfaces"
)

func (p *ScenePlayer) GetSystem(name string) systems.ISystem {
	return p.systems.GetSystem(name)
}

// 注册system
// . 在scene_systems_visit中添加模块访问
// . 在下面明确的创建

func (p *ScenePlayer) createSystems() {
	p.systems = systems.NewSystems(p)

	// create all systems
	sysfactory.GetFactory().Visit(func(name string, creator factory.ICreator) {
		p.createSystem(name)
	})

	p.GetSystem("test2").(interfaces.Test2).Hello()
}

func (p *ScenePlayer) createSystem(name string) {
	obj := sysfactory.GetFactory().Create(name)
	if obj == nil {
		return
	}
	sys, ok := obj.(systems.ISystem)
	if !ok || sys == nil {
		l.L.Errorf("error system create: %v", name)
		return
	}

	p.systems.AddSystem(name, sys)
}

func (p *ScenePlayer) systemsLoadData(info *mymsg.PlayerInfo) {
	p.systems.Visit(func(sys systems.ISystem) {
		sys.LoadData(info)
	})
}

func (p *ScenePlayer) systemsInitData() {
	p.systems.Visit(func(sys systems.ISystem) {
		sys.InitData()
	})
}

func (p *ScenePlayer) systemsSaveData(info *mymsg.PlayerInfo) {
	p.systems.Visit(func(sys systems.ISystem) {
		sys.SaveData(info)
	})
}

func (p *ScenePlayer) systemsOnEnterWorld(switchLine bool) {
	p.systems.Visit(func(sys systems.ISystem) {
		sys.OnEnterWorld(switchLine)
	})
}

func (p *ScenePlayer) systemsReqInfo() {
	p.systems.Visit(func(sys systems.ISystem) {
		sys.PushInfoToClient()
	})
}

func (p *ScenePlayer) SystemRequest(system, cmd string, args []byte, cb func(ret []byte, errCode int32)) {
	p.systems.Request(system, cmd, args, cb)
}
