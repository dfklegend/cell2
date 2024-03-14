package module

import (
	"sync"

	"github.com/dfklegend/cell2/baseapp/interfaces"
	"github.com/dfklegend/cell2/utils/common"
	l "github.com/dfklegend/cell2/utils/logger"
)

//	Module列表
type ModList struct {
	lock *sync.RWMutex
	mods []interfaces.IAppModule
}

func NewModList() *ModList {
	return &ModList{
		lock: &sync.RWMutex{},
		mods: make([]interfaces.IAppModule, 0),
	}
}

func (m *ModList) AddModule(mod interfaces.IAppModule) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.mods = append(m.mods, mod)
}

func (m *ModList) Filter(startOrder bool, doFunc func(mod interfaces.IAppModule, next interfaces.FuncWithSucc),
	finish interfaces.FuncWithSucc) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if len(m.mods) == 0 {
		finish(true)
		return
	}

	var doNow func()
	var next func(bool)

	if startOrder {
		// 正向
		index := 0
		next = func(succ bool) {
			if !succ {
				finish(false)
				return
			}
			index += 1
			doNow()
		}

		doNow = func() {
			if index >= len(m.mods) {
				finish(true)
				return
			}

			cur := m.mods[index]
			doFunc(cur, next)
		}
	} else {
		// 反向
		index := len(m.mods) - 1
		next = func(succ bool) {
			if !succ {
				finish(false)
				return
			}

			index -= 1
			doNow()
		}

		doNow = func() {
			if index < 0 {
				finish(true)
				return
			}

			cur := m.mods[index]
			doFunc(cur, next)
		}
	}

	doNow()
}

func (m *ModList) Start(finish interfaces.FuncWithSucc) {
	m.Filter(true, func(mod interfaces.IAppModule, next interfaces.FuncWithSucc) {
		defer func() {
			if err := recover(); err != nil {
				l.E.Errorf("panic in module.Start:%v", err)
				stack := common.GetStackStr()
				l.E.Infof(stack)
			}
		}()

		mod.Start(next)
	}, finish)
}

//	反向

func (m *ModList) Stop(finish interfaces.FuncWithSucc) {
	m.Filter(false, func(mod interfaces.IAppModule, next interfaces.FuncWithSucc) {
		defer func() {
			if err := recover(); err != nil {
				l.E.Errorf("panic in module.Stop:%v", err)
				stack := common.GetStackStr()
				l.E.Infof(stack)
			}
		}()

		mod.Stop(next)
	}, finish)
}
