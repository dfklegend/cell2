package lua

import (
	"fmt"
	"log"
	"testing"
	"unsafe"

	"github.com/dfklegend/cell2/utils/event/light"
	"github.com/dfklegend/cell2/utils/golua"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

func Test_Normal(t *testing.T) {
	golua.InitLuaPathAndCompile("./servicetestscripts", true)

	goEnv := NewEnv()
	s := NewBuilder().
		Prepare().
		BindUserTypes(func(L *lua.LState) {
			BindTestScriptEnv(L)
		}).Start(goEnv, "main.lua", "start").
		GetService()
	e := s.engine

	fmt.Printf("%v %v %v\n", goEnv.V, goEnv.F32, goEnv.F64)

	global := NewFuncs()
	global.AddFunc(e.L.Env, "createSkillScript")
	global.AddFunc(e.L.Env, "testEnv")

	golua.Call(e.L, global.Find("testEnv"), goEnv)
	fmt.Printf("%v %v %v\n", goEnv.V, goEnv.F32, goEnv.F64)

	skillApis := NewFuncs()
	skillApis.AddFunc(GetTable(e.L.Env, "Root.Game.skillAPIs"), "skill_onSkillHit")

	ret, err := golua.CallWithResult(e.L, global.Find("createSkillScript"), "example")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret)

	golua.Call(e.L, skillApis.Find("skill_onSkillHit"), ret, 1)
	golua.Call(e.L, skillApis.Find("skill_onSkillHit"), ret, 1)
	golua.Call(e.L, skillApis.Find("skill_onSkillHit"), ret, 1)
	golua.Call(e.L, skillApis.Find("skill_onSkillHit"), ret, 1)
	s.Stop()
}

type TestData struct {
	Str string
}

func BindTestData(L *lua.LState) {
	L.SetGlobal("TestData", luar.NewType(L, TestData{}))
}

func Test_UserData(t *testing.T) {
	golua.InitLuaPathAndCompile("./servicetestscripts", true)

	s := NewBuilder().
		Prepare().
		BindUserTypes(func(L *lua.LState) {
			BindTestScriptEnv(L)
			BindTestData(L)
		}).Start(NewEnv(), "main.lua", "start").
		GetService()
	e := s.engine

	global := NewFuncs()
	global.AddFunc(e.L.Env, "createSkillScript")

	skillApis := NewFuncs()
	skillApis.AddFunc(GetTable(e.L.Env, "Root.Game.skillAPIs"), "skill_onSkillHit")
	skillApis.AddFunc(GetTable(e.L.Env, "Root.Game.skillAPIs"), "skill_test")

	data := &TestData{
		Str: "some",
	}

	golua.Call(e.L, skillApis.Find("skill_test"), data)
	s.Stop()
}

// 测试Lua GC
//
func Test_GC(t *testing.T) {
	golua.InitLuaPathAndCompile("./servicetestscripts", true)

	s := NewBuilder().
		Prepare().
		BindUserTypes(func(L *lua.LState) {
			BindTestScriptEnv(L)
		}).Start(NewEnv(), "main.lua", "start").
		GetService()
	e := s.engine

	global := NewFuncs()
	global.AddFunc(e.L.Env, "createSkillScript")

	skillApis := NewFuncs()
	skillApis.AddFunc(GetTable(e.L.Env, "Root.Game.skillAPIs"), "skill_onSkillHit")

	ret, err := golua.CallWithResult(e.L, global.Find("createSkillScript"), "example")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret)

	address := uintptr(unsafe.Pointer(ret.(*lua.LTable)))
	golua.Call(e.L, skillApis.Find("skill_onSkillHit"), ret, 1)

	// 注: 如果将ret赋值成nil，那么ret对应的table将在后续被GC掉
	// golua.Call(e.L, skillApis.Find("skill_onSkillHit"), ret1, 1)将会异常
	// -- A --
	//ret = nil
	for i := 0; i < 3; i++ {
		e.DoLuaString(`collectgarbage()`)
	}

	// -- B --
	ret1 := (*lua.LTable)(unsafe.Pointer(address))
	golua.Call(e.L, skillApis.Find("skill_onSkillHit"), ret1, 1)

	// 如果上面A处关闭，但是没有下面C这个使用
	// 上面B处的Call也会出错(随机出错)，怀疑是优化导致GC还是释放掉了(后面没引用)
	// -- C --
	fmt.Println(ret)
	s.Stop()
}

// 测试事件，参数传递
func Test_Events(t *testing.T) {
	golua.InitLuaPathAndCompile("./servicetestscripts", true)

	goEnv := NewEnv()
	s := NewBuilder().
		Prepare().
		BindUserTypes(func(L *lua.LState) {
			BindTestScriptEnv(L)
			L.SetGlobal("EventCenter", luar.NewType(L, light.EventCenter{}))
		}).Start(goEnv, "main.lua", "start").
		GetService()
	e := s.engine

	global := NewFuncs()
	global.AddFunc(e.L.Env, "testInitEvents")

	events := light.NewEventCenter()

	golua.Call(e.L, global.Find("testInitEvents"), events)

	events.Publish("hello", goEnv)
	fmt.Printf("%v %v %v\n", goEnv.V, goEnv.F32, goEnv.F64)
}

// 测试用户变量传递
// 复杂struct统一转成LUserData，不管有没有注册类型
// 注册类型只是为了方便脚本创建
func Test_UserData1(t *testing.T) {
	golua.InitLuaPathAndCompile("./servicetestscripts", true)

	goEnv := NewEnv()
	s := NewBuilder().
		Prepare().
		BindUserTypes(func(L *lua.LState) {
			BindTestScriptEnv(L)
			BindTestUserData(L)
			L.SetGlobal("EventCenter", luar.NewType(L, light.EventCenter{}))
		}).Start(goEnv, "main.lua", "start").
		GetService()
	e := s.engine

	global := NewFuncs()
	global.AddFunc(e.L.Env, "testUserData")

	data := NewUserData()
	data.Data1 = &TestUserData1{}

	golua.Call(e.L, global.Find("testUserData"), data)

	fmt.Printf("%v %v\n", data.V, data.Data1.V)
}

// 测试
func Benchmark_Normal(b *testing.B) {
	golua.InitLuaPathAndCompile("./servicetestscripts", true)

	s := NewBuilder().
		Prepare().
		BindUserTypes(func(L *lua.LState) {
			BindTestScriptEnv(L)
		}).Start(NewEnv(), "main.lua", "start").
		GetService()
	e := s.engine

	global := NewFuncs()
	global.AddFunc(e.L.Env, "createSkillScript")

	skillApis := NewFuncs()
	skillApis.AddFunc(GetTable(e.L.Env, "Root.Game.skillAPIs"), "skill_onSkillHit")

	//ret, err := e.DoLuaMethodWithResult("skills/mgr.lua", "createSkillScript", "example")
	ret, err := golua.CallWithResult(e.L, global.Find("createSkillScript"), "example")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ret)

	b.ResetTimer()

	log.Println("bench start")
	for i := 0; i < b.N; i++ {
		golua.Call(e.L, skillApis.Find("skill_onSkillHit"), ret, 1)
	}
	s.Stop()
}
