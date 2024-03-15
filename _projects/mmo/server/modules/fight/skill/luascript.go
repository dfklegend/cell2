package skill

import (
	"fmt"

	"github.com/dfklegend/cell2/utils/golua"
	l "github.com/dfklegend/cell2/utils/logger"
	lua "github.com/yuin/gopher-lua"

	lua2 "mmo/modules/lua"

	"mmo/modules/fight/common"
	"mmo/modules/fight/lua/env"
)

// ScriptMgr for skill script
// 每个脚本环境唯一
// 挂在ScriptEnvData上
// 技能脚本都是单体对象
// 在对应技能释放时，可以触发一些回调，可以做逻辑定制
type ScriptMgr struct {
	env *env.ScriptEnvData
	// value can be nil
	scripts map[string]*Script
}

func newScriptMgr(env *env.ScriptEnvData) *ScriptMgr {
	return &ScriptMgr{
		env:     env,
		scripts: map[string]*Script{},
	}
}

func (m *ScriptMgr) addScript(name string, script *Script) {
	m.scripts[name] = script
}

func (m *ScriptMgr) GetScript(name string) *Script {
	script, ok := m.scripts[name]
	if ok {
		return script
	}
	obj := m.tryCreateScript(name)
	if obj == nil {
		m.addScript(name, nil)
		return nil
	}
	script = newScript()
	script.Init(m.env.State, obj)
	m.addScript(name, script)
	return script
}

func (m *ScriptMgr) tryCreateScript(name string) *lua.LTable {
	global := m.env.GlobalAPIs
	L := m.env.State
	ret, err := golua.CallWithResult(L, global.Find("createSkillScript"), name)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if ret == lua.LNil {
		return nil
	}

	if ret.Type() != lua.LTTable {
		return nil
	}
	return ret.(*lua.LTable)
}

type Script struct {
	state *lua.LState
	obj   *lua.LTable

	fnChangeSkill lua.LValue
	fnStart       lua.LValue
	fnHit         lua.LValue
}

func newScript() *Script {
	return &Script{}
}

func (s *Script) Init(L *lua.LState, obj *lua.LTable) {
	s.state = L
	s.obj = obj

	s.fnChangeSkill = lua2.GetObjFunc(L, obj, "changeSkill")
	s.fnStart = lua2.GetObjFunc(L, obj, "onSkillStart")
	s.fnHit = lua2.GetObjFunc(L, obj, "onSkillHit")
}

func (s *Script) onSkillStart(proxy *Proxy) {
	if s.fnStart == nil {
		return
	}
	golua.Call(s.state, s.fnStart, s.obj, proxy)
}

func (s *Script) onSkillHit(proxy *Proxy) {
	if s.fnHit == nil {
		return
	}
	golua.Call(s.state, s.fnHit, s.obj, proxy)
}

func (s *Script) changeSkill(src common.ICharacter) common.SkillId {
	if s.fnChangeSkill == nil {
		return ""
	}
	ret, err := golua.CallWithResult(s.state, s.fnChangeSkill, s.obj, src.GetProxy())
	if err != nil {
		l.L.Errorf("lua changeskill error: %v", err)
		return ""
	}
	if ret == lua.LNil {
		return ""
	}

	if ret.Type() != lua.LTString {
		return ""
	}
	return ret.String()
}

func CreateSkillScriptMgr(env *env.ScriptEnvData) {
	env.SetValue("skillScripts", newScriptMgr(env))
}
