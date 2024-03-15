package skill

import (
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"

	"mmo/modules/fight/common"
)

// Proxy lua proxy for skill
type Proxy struct {
	skill *Skill
}

func newProxy(skill *Skill) *Proxy {
	return &Proxy{
		skill: skill,
	}
}

func (p *Proxy) GetId() string {
	return p.skill.id
}

func (p *Proxy) IsBGSkill() bool {
	return p.skill.bgSkill
}

func (p *Proxy) Owner() common.ICharProxy {
	return p.skill.owner.GetProxy()
}

func (p *Proxy) Tar() common.ICharProxy {
	tar := p.skill.world.GetChar(p.skill.tar)
	if tar == nil {
		return nil
	}
	return tar.GetProxy()
}

func BindSkill(L *lua.LState) {
	L.SetGlobal("Skill", luar.NewType(L, Proxy{}))
}
