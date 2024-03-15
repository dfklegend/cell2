package buf

import (
	"mmo/modules/fight/script"
)

type FnCreator func() script.IBufScript

var mgr = newMgr()

type bufMgr struct {
	creators map[string]FnCreator
}

func newMgr() *bufMgr {
	return &bufMgr{
		creators: map[string]FnCreator{},
	}
}

func (m *bufMgr) Register(name string, fn FnCreator) {
	m.creators[name] = fn
}

func (m *bufMgr) Create(name string) script.IBufScript {
	fn := m.creators[name]
	if fn == nil {
		return nil
	}
	return fn()
}
