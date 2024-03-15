package luabind

import (
	"github.com/dfklegend/cell2/utils/event/light"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

func BindEventCenter(L *lua.LState) {
	L.SetGlobal("EventCenter", luar.NewType(L, light.EventCenter{}))
}
