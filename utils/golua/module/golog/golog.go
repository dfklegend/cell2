package golog

import (
	"fmt"
	"log"

	"github.com/yuin/gopher-lua"
)

func Loader(L *lua.LState) int {
	// register functions to the table
	mod := L.SetFuncs(L.NewTable(), exports)
	// register other stuff
	//L.SetField(mod, "name", lua.LString("value"))
	// returns the module
	L.Push(mod)
	return 1
}

func Preload(L *lua.LState) {
	L.PreloadModule("go_log", Loader)
}

var exports = map[string]lua.LGFunction{
	"debug": debug,
	"info":  info,
	"warn":  warn,
	"error": error,
}

func debug(L *lua.LState) int {
	log.Println(fmt.Sprintf("[lua][debug]%s", L.CheckString(-1)))
	return 0
}

func info(L *lua.LState) int {
	log.Println(fmt.Sprintf("[lua][info]%s", L.CheckString(-1)))
	return 0
}

func warn(L *lua.LState) int {
	log.Println(fmt.Sprintf("[lua][warn]%s", L.CheckString(-1)))
	return 0
}

func error(L *lua.LState) int {
	log.Println(fmt.Sprintf("[lua][error]%s", L.CheckString(-1)))
	return 0
}
