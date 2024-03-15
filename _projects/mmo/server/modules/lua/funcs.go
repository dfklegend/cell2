package lua

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// Funcs
// 将Lua代码中的函数统一查找出来
// 便于使用
type Funcs struct {
	funcs map[string]lua.LValue
}

func NewFuncs() *Funcs {
	return &Funcs{
		funcs: map[string]lua.LValue{},
	}
}

func (f *Funcs) AddFunc(table *lua.LTable, name string) {
	fn := table.RawGetString(name)
	if fn == nil || fn.Type() != lua.LTFunction {
		fmt.Printf("can not found lua method! %s\n", name)
		return
	}
	f.funcs[name] = fn
}

func (f *Funcs) Find(name string) lua.LValue {
	fn, ok := f.funcs[name]
	if !ok {
		return nil
	}
	return fn
}
