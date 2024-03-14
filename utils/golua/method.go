package golua

import (
	"errors"

	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

var errorNilFunction = errors.New("nil function")

func IsLuaFunction(fn lua.LValue) bool {
	if fn == nil || fn.Type() != lua.LTFunction {
		return false
	}
	return true
}

func Call(state *lua.LState, fn lua.LValue, args ...any) error {
	if fn == nil {
		panic(errorNilFunction)
	}
	var argsArr []lua.LValue

	for _, arg := range args {
		argsArr = append(argsArr, luar.New(state, arg))
	}

	err := state.CallByParam(lua.P{
		Fn:      fn,
		NRet:    0,
		Protect: true,
		Handler: nil,
	}, argsArr...)

	return checkError(err)
}

func CallWithResult(state *lua.LState, fn lua.LValue, args ...any) (lua.LValue, error) {
	if fn == nil {
		panic(errorNilFunction)
	}

	var argsArr []lua.LValue

	for _, arg := range args {
		argsArr = append(argsArr, luar.New(state, arg))
	}

	err := state.CallByParam(lua.P{
		Fn:      fn,
		NRet:    1,
		Protect: true,
		Handler: nil,
	}, argsArr...)

	if err != nil {
		checkError(err)
	}

	ret := state.Get(-1)
	state.Pop(1)

	return ret, err
}
