package main

import (
	"fmt"
	"os"

	lua "github.com/yuin/gopher-lua"

	"github.com/dfklegend/cell2/utils/golua"
)

func main() {

	fmt.Println()
	//设置lua目录
	//golua.SetUserLuaPath("exception", false)
	golua.SetUserLuaPath("luaScript", true)
	fmt.Printf("%s=%s", lua.LuaPath, os.Getenv(lua.LuaPath))
	fmt.Println()

	luaEngine := golua.NewLuaEngine()
	defer luaEngine.Close()

	luaEngine.DoLuaString("print('DoLuaString')")

	err := luaEngine.DoLuaFile("init.lua")
	if err != nil {
		panic(err)
	}

	luaEngine.DoLuaMethod("test.lua", "sayhello")
	luaEngine.DoLuaMethod("aa/test2.lua", "sayhello")
	//luaEngine.DoLuaMethodWithResult("example.lua", "hello")
	//luaEngine.DoLuaFile("example.lua")
	//luaEngine.DoLuaMethodWithResult("example.lua", "triggerError")
}
