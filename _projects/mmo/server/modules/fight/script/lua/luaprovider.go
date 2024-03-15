package lua

import (
	"fmt"

	"github.com/dfklegend/cell2/utils/golua"
	lua "github.com/yuin/gopher-lua"

	"mmo/common/config"
	"mmo/modules/fight/lua/env"
	"mmo/modules/fight/script"
)

type ScriptProvider struct {
	env *env.ScriptEnvData
}

func newProvider(env *env.ScriptEnvData) *ScriptProvider {
	return &ScriptProvider{
		env: env,
	}
}

func (p *ScriptProvider) CreateBufScript(name string) script.IBufScript {
	if config.EnableFightLuaPlugin == 0 {
		return nil
	}

	env := p.env

	global := env.GlobalAPIs
	L := env.State
	ret, err := golua.CallWithResult(L, global.Find("createBufScript"), name)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if ret.Type() != lua.LTTable {
		return nil
	}

	obj := ret.(*lua.LTable)
	script := newBufScript()
	script.Init(L, obj)
	return script
}
