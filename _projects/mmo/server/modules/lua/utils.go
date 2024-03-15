package lua

import (
	"strings"

	"github.com/dfklegend/cell2/utils/golua"
	lua "github.com/yuin/gopher-lua"
)

// GetTable
// 从table获取嵌套的table
// 比如GetTable(L.Env, "Root.Game.skillAPIs")
func GetTable(table *lua.LTable, path string) *lua.LTable {
	index := strings.Index(path, ".")
	head := ""
	tail := ""
	if index == -1 {
		head = path
	} else {
		head = path[:index]
		tail = path[index+1:]
	}
	sub := table.RawGetString(head)
	if sub == nil || sub.Type() != lua.LTTable {
		return nil
	}
	if index == -1 {
		return sub.(*lua.LTable)
	}
	return GetTable(sub.(*lua.LTable), tail)
}

func GetObjFunc(L *lua.LState, obj *lua.LTable, name string) lua.LValue {
	fn := L.GetField(obj, name)
	if fn == nil || !golua.IsLuaFunction(fn) {
		return nil
	}
	return fn
}
