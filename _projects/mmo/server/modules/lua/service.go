package lua

import (
	"log"

	"github.com/dfklegend/cell2/utils/golua"
	lua "github.com/yuin/gopher-lua"
)

/*
	gopher-lua的GC交给go来处理
	Lua部分相当于都是使用go数据结构来执行
	所以只要go部分持有对象，对象就有效
*/

// Service
// 环境和启动不一样
type Service struct {
	engine  *golua.LuaEngine
	goEnv   IGoEnv
	envData IScriptEnvData
}

func NewService() *Service {
	return &Service{
		engine: golua.NewLuaEngine(),
	}
}

func (s *Service) GetEngine() *golua.LuaEngine {
	return s.engine
}

func (s *Service) GetL() *lua.LState {
	return s.engine.L
}

func (s *Service) Prepare() {
	s.engine.LoadGopherLuaLibs()
}

func (s *Service) Start(goEnv IGoEnv, luaFile string, funcName string) {
	s.goEnv = goEnv
	if err := s.engine.DoLuaMethod(luaFile, funcName, goEnv); err != nil {
		log.Println(err)
	}
}

func (s *Service) SetEnvData(data IScriptEnvData) {
	s.envData = data
}

func (s *Service) GetEnvData() IScriptEnvData {
	return s.envData
}

func (s *Service) Stop() {
	s.engine.Close()
}
