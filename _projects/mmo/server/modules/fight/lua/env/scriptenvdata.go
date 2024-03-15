package env

import (
	lua "github.com/yuin/gopher-lua"

	lua2 "mmo/modules/lua"
)

type ScriptEnvData struct {
	State *lua.LState
	// 全局函数
	GlobalAPIs *lua2.Funcs
	SkillAPIs  *lua2.Funcs

	values map[string]any
}

func NewScriptEnvData() *ScriptEnvData {
	return &ScriptEnvData{
		GlobalAPIs: lua2.NewFuncs(),
		SkillAPIs:  lua2.NewFuncs(),
		values:     map[string]any{},
	}
}

func (m *ScriptEnvData) Prepare(state *lua.LState) {
	m.State = state

	m.GlobalAPIs.AddFunc(state.Env, "createSkillScript")
	m.GlobalAPIs.AddFunc(state.Env, "createBufScript")
	m.GlobalAPIs.AddFunc(state.Env, "globalhello")

	t := lua2.GetTable(state.Env, "Root.Game.skillAPIs")
	m.SkillAPIs.AddFunc(t, "skill_onSkillHit")
	m.SkillAPIs.AddFunc(t, "buf_onAdd")
}

func (m *ScriptEnvData) SetValue(key string, v any) {
	m.values[key] = v
}

func (m *ScriptEnvData) GetValue(key string) any {
	return m.values[key]
}
