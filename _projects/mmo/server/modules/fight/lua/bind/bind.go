package luabind

import (
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"

	"mmo/modules/fight/buf"
	charimpls "mmo/modules/fight/character/impls"
	"mmo/modules/fight/common"
	"mmo/modules/fight/lua/env"
	"mmo/modules/fight/skill"
)

// BindAllFightTypes 绑定所有的战斗相关类型
func BindAllFightTypes(L *lua.LState) {
	BindDmgInstance(L)
	skill.BindSkill(L)
	buf.BindBuf(L)
	charimpls.BindCharacter(L)
}

// InitFightGoAPIs 初始化一些api给lua
func InitFightGoAPIs(L *lua.LState) {

}

func MakeScriptEnvData(L *lua.LState) *env.ScriptEnvData {
	d := env.NewScriptEnvData()
	d.Prepare(L)
	buf.CreateBufScriptMgr(d)
	skill.CreateSkillScriptMgr(d)
	return d
}

func BindDmgInstance(L *lua.LState) {
	L.SetGlobal("DmgInstance", luar.NewType(L, common.DmgInstance{}))
}
