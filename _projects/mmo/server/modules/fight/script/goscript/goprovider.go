package goscript

import (
	"mmo/modules/fight/script"
)

type Provider struct {
	bufMgr script.IBufScriptMgr
}

func newProvider(bufMgr script.IBufScriptMgr) *Provider {
	return &Provider{
		bufMgr: bufMgr,
	}
}

func (p *Provider) CreateBufScript(name string) script.IBufScript {
	return p.bufMgr.Create(name)
}
