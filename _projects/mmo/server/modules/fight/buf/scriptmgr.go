package buf

import (
	"mmo/modules/fight/lua/env"
)

// ScriptMgr 记录不存在的脚本
// 减少无效lua调用
type ScriptMgr struct {
	env          *env.ScriptEnvData
	emptyScripts map[string]string
}

func newMgr(env *env.ScriptEnvData) *ScriptMgr {
	return &ScriptMgr{
		env:          env,
		emptyScripts: map[string]string{},
	}
}

func (m *ScriptMgr) isEmptyScript(name string) bool {
	_, ok := m.emptyScripts[name]
	if ok {
		return true
	}
	return false
}

func (m *ScriptMgr) addEmptyScript(name string) {
	m.emptyScripts[name] = name
}

func CreateBufScriptMgr(env *env.ScriptEnvData) {
	env.SetValue("bufScripts", newMgr(env))
}
