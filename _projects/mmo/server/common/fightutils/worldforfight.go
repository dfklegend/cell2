package fightutils

import (
	"mmo/common/entity"
	"mmo/modules/fight/common"
	"mmo/modules/lua"
	"mmo/servers/scene/entity/components"
	define2 "mmo/servers/scene/entity/define"
)

//

type WorldForFight struct {
	world        entity.IWorld
	timeProvider common.ITimeProvider
	watcher      common.IWatcher
	recorder     common.IFightDetailRecorder
	lua          *lua.Service
	scriptMgr    common.IScriptMgr
}

func NewWorldForFight(world entity.IWorld, provider common.ITimeProvider,
	watcher common.IWatcher, recorder common.IFightDetailRecorder,
	lua *lua.Service, scriptMgr common.IScriptMgr) *WorldForFight {
	return &WorldForFight{
		world:        world,
		timeProvider: provider,
		watcher:      watcher,
		recorder:     recorder,
		lua:          lua,
		scriptMgr:    scriptMgr,
	}
}

func (w *WorldForFight) GetChar(id common.CharId) common.ICharacter {
	tar := w.world.GetEntity(id)
	if tar == nil {
		return nil
	}
	c := tar.GetComponent(define2.BaseUnit)
	if c == nil {
		return nil
	}
	u := c.(*components.BaseUnit)
	if u == nil {
		return nil
	}
	return u.GetChar()
}

func (w *WorldForFight) GetTimeProvider() common.ITimeProvider {
	return w.timeProvider
}

func (w *WorldForFight) GetWatcher() common.IWatcher {
	return w.watcher
}

func (w *WorldForFight) GetDetailRecorder() common.IFightDetailRecorder {
	return w.recorder
}

func (w *WorldForFight) GetLua() *lua.Service {
	return w.lua
}

func (w *WorldForFight) GetScriptMgr() common.IScriptMgr {
	return w.scriptMgr
}
