package sceneservice

import (
	"github.com/dfklegend/cell2/node/service"

	luabind2 "mmo/common/luabind"
	"mmo/modules/fight/lua/bind"

	lua2 "mmo/modules/lua"

	lua "github.com/yuin/gopher-lua"
)

func initLua() *lua2.Service {
	s := lua2.NewBuilder().
		Prepare().
		PrepareNodes("nodeinit.lua").
		BindUserTypes(func(L *lua.LState) {
			luabind2.BindEventCenter(L)
			luabind.BindAllFightTypes(L)

		}).
		PreNext(func(L *lua.LState) {
			luabind.InitFightGoAPIs(L)
		}).
		Start(nil, "main.lua", "start").
		GetService()
	envData := luabind.MakeScriptEnvData(s.GetL())
	s.SetEnvData(envData)
	return s
}

type LuaComponent struct {
	*service.BaseComponent
	lua *lua2.Service
}

func NewLuaComponent(lua *lua2.Service) *LuaComponent {
	return &LuaComponent{
		BaseComponent: service.NewBaseComponent(),
		lua:           lua,
	}
}

func (c *LuaComponent) OnAdd() {
}

func (c *LuaComponent) OnRemove() {
}
