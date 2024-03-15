package lua

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

type Builder struct {
	service *Service
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Prepare() *Builder {
	b.service = NewService()
	b.service.Prepare()
	return b
}

// PrepareNodes 提前初始化节点
func (b *Builder) PrepareNodes(luaFile string) *Builder {
	if err := b.service.GetEngine().DoLuaFile(luaFile); err != nil {
		log.Println(err)
	}
	return b
}

// BindUserTypes 做一些用户类型绑定
func (b *Builder) BindUserTypes(doFunc func(L *lua.LState)) *Builder {
	doFunc(b.service.GetL())
	return b
}

// PrepareUserModules 准备用户模块
func (b *Builder) PrepareUserModules(doFunc func(L *lua.LState)) *Builder {
	doFunc(b.service.GetL())
	return b
}

// PreNext 一些其他初始化
func (b *Builder) PreNext(doFunc func(L *lua.LState)) *Builder {
	doFunc(b.service.GetL())
	return b
}

// Start 启动执行
func (b *Builder) Start(env IGoEnv, luaFile string, funcName string) *Builder {
	b.service.Start(env, luaFile, funcName)
	return b
}

// PostNext 启动后一些初始化
func (b *Builder) PostNext(doFunc func(L *lua.LState)) *Builder {
	doFunc(b.service.GetL())
	return b
}

func (b *Builder) GetService() *Service {
	return b.service
}
