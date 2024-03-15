package mgr

import (
	"mmo/modules/fight/script"
)

type Mgr struct {
	providers []script.IScriptProvider
}

func NewMgr() *Mgr {
	return &Mgr{
		providers: make([]script.IScriptProvider, 0),
	}
}

func (m *Mgr) AddProvider(provider script.IScriptProvider) {
	m.providers = append(m.providers, provider)
}

func (m *Mgr) CreateBufScript(name string) script.IBufScript {
	var script script.IBufScript
	for i := 0; i < len(m.providers); i++ {
		script = m.providers[i].CreateBufScript(name)
		if script != nil {
			return script
		}
	}
	return nil
}
