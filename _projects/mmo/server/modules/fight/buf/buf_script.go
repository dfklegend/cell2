package buf

import (
	"fmt"

	"github.com/dfklegend/cell2/utils/golua"
	lua "github.com/yuin/gopher-lua"

	"mmo/common/config"
	"mmo/modules/fight/common"
	"mmo/modules/fight/lua/env"
	"mmo/modules/fight/script"
)

// 优化: 如果某个为nil,则肯定不存在脚本
func (b *Buf) tryCreateScript() bool {
	if config.EnableFightLuaPlugin == 0 {
		return false
	}
	var m0 common.IScriptMgr
	if m0 = b.owner.GetWorld().GetScriptMgr(); m0 == nil {
		return false
	}

	mgr := m0.(script.IScriptMgr)
	b.script = mgr.CreateBufScript(b.id)
	return true
}

func (b *Buf) oldTryCreateScript() bool {
	if config.EnableFightLuaPlugin == 0 {
		return false
	}
	if b.owner.GetWorld().GetLua() == nil {
		return false
	}
	env := b.owner.GetWorld().GetLua().GetEnvData().(*env.ScriptEnvData)
	mgr := env.GetValue("bufScripts").(*ScriptMgr)
	if mgr.isEmptyScript(b.id) {
		return false
	}
	global := env.GlobalAPIs
	L := env.State
	ret, err := golua.CallWithResult(L, global.Find("createBufScript"), b.id)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if ret == lua.LNil {
		mgr.addEmptyScript(b.id)
		return false
	}

	if ret.Type() != lua.LTTable {
		return false
	}

	//obj := ret.(*lua.LTable)
	//b.script = newScript()
	//b.script.Init(L, obj)
	return true
}
