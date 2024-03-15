package lua

import (
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

type TestUserData1 struct {
	V   int
	F32 float32
	F64 float64
}

type TestUserData struct {
	Data1 *TestUserData1
	V     int
	F32   float32
	F64   float64
}

func NewUserData() *TestUserData {
	return &TestUserData{}
}

func BindTestUserData(L *lua.LState) {
	L.SetGlobal("testUserData", luar.NewType(L, TestUserData{}))
}
