package lua

import (
	"github.com/dfklegend/cell2/utils/golua"
	lua "github.com/yuin/gopher-lua"

	"mmo/modules/fight/script"
	lua2 "mmo/modules/lua"
)

type BufScript struct {
	state *lua.LState
	obj   *lua.LTable

	fnOnStart   lua.LValue
	fnOnEnd     lua.LValue
	fnOnTriggle lua.LValue
}

func newBufScript() *BufScript {
	return &BufScript{}
}

func (s *BufScript) Init(args ...any) {
	L := args[0].(*lua.LState)
	obj := args[1].(*lua.LTable)
	s.doInit(L, obj)
}

func (s *BufScript) doInit(L *lua.LState, obj *lua.LTable) {
	s.state = L
	s.obj = obj

	s.fnOnStart = lua2.GetObjFunc(L, obj, "onStart")
	s.fnOnEnd = lua2.GetObjFunc(L, obj, "onEnd")
	s.fnOnTriggle = lua2.GetObjFunc(L, obj, "onTriggle")
}

func (s *BufScript) OnStart(proxy script.IBufProxy) {
	if s.fnOnStart == nil {
		return
	}
	golua.Call(s.state, s.fnOnStart, s.obj, proxy)
}

func (s *BufScript) OnEnd(proxy script.IBufProxy) {
	if s.fnOnEnd == nil {
		return
	}
	golua.Call(s.state, s.fnOnEnd, s.obj, proxy)
}

func (s *BufScript) OnTriggle(proxy script.IBufProxy) {
	if s.fnOnTriggle == nil {
		return
	}
	golua.Call(s.state, s.fnOnTriggle, s.obj, proxy)
}
