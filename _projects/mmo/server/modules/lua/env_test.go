package lua

import (
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

// TestScriptEnv
// for test
type TestScriptEnv struct {
	V   int
	F32 float32
	F64 float64
}

func NewEnv() *TestScriptEnv {
	return &TestScriptEnv{}
}

func BindTestScriptEnv(L *lua.LState) {
	L.SetGlobal("testScriptEnv", luar.NewType(L, TestScriptEnv{}))
}
