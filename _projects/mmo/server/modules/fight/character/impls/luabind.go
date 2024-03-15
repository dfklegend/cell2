package charimpls

import (
	"github.com/dfklegend/cell2/utils/event/light"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"

	"mmo/modules/fight/attr"
	"mmo/modules/fight/common"
	env2 "mmo/modules/fight/lua/env"
)

// Proxy proxy for character
type Proxy struct {
	c *Character
	// 用户数据
	// 脚本层使用
	UserData *lua.LTable
}

func newProxy(c *Character) *Proxy {
	if c.world.GetLua() == nil {
		return nil
	}
	env := c.world.GetLua().GetEnvData().(*env2.ScriptEnvData)
	return &Proxy{
		c:        c,
		UserData: env.State.NewTable(),
	}
}

func (p *Proxy) GetId() int32 {
	return int32(p.c.id)
}

func (p *Proxy) GetEvents() *light.EventCenter {
	return p.c.events
}

func (p *Proxy) OffsetBase(index int, value attr.Value) {
	p.c.OffsetBase(index, value)
}

func (p *Proxy) IsDead() bool {
	return p.c.IsDead()
}

func (p *Proxy) GetTarId() common.CharId {
	return p.c.GetTarId()
}

func (p *Proxy) CallbackSkill(id common.SkillId, level int, src *Proxy, tar common.CharId) {
	p.c.CallbackSkill(id, level, src.c, tar)
}

func (p *Proxy) GoCallbackSkill(id common.SkillId, level int, src common.ICharProxy, tar common.CharId) {
	sp := src.(*Proxy)
	p.c.CallbackSkill(id, level, sp.c, tar)
}

func (p *Proxy) AddBuf(id common.BufId, level int, stack int) {
	p.c.AddBuf(p.c, p.c, id, level, stack)
}

func (p *Proxy) RemoveBuf(id common.BufId) {
	p.c.RemoveBuf(id)
}

func BindCharacter(L *lua.LState) {
	L.SetGlobal("Character", luar.NewType(L, Proxy{}))
}
