package test

import (
	"mmo/modules/fight/common"
	"mmo/modules/lua"
)

type TestWorld struct {
	chars        []common.ICharacter
	timeProvider common.ITimeProvider
	watcher      common.IWatcher
}

func newTestWorld(timeProvider common.ITimeProvider, watcher common.IWatcher) *TestWorld {
	return &TestWorld{
		chars:        []common.ICharacter{},
		timeProvider: timeProvider,
		watcher:      watcher,
	}
}

func (w *TestWorld) GetTimeProvider() common.ITimeProvider {
	return w.timeProvider
}

func (w *TestWorld) GetWatcher() common.IWatcher {
	return w.watcher
}

func (w *TestWorld) GetDetailRecorder() common.IFightDetailRecorder {
	return nil
}

func (w *TestWorld) AddChar(c common.ICharacter) {
	w.chars = append(w.chars, c)
}

func (w *TestWorld) GetChar(id common.CharId) common.ICharacter {
	for _, v := range w.chars {
		if v.GetId() == id {
			return v
		}
	}
	return nil
}

func (w *TestWorld) GetLua() *lua.Service {
	return nil
}

func (w *TestWorld) GetScriptMgr() common.IScriptMgr {
	return nil
}
