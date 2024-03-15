package script

import (
	"mmo/modules/fight/common"
)

type ISkillProxy interface {
	GetId() string
	IsBGSkill() bool
	Owner() common.ICharProxy
	Tar() common.ICharProxy
}

// IBufProxy
// 脚本环境里代表buf
type IBufProxy interface {
	GetId() common.BufId
	Owner() common.ICharProxy
}

type IBufScript interface {
	Init(args ...any)
	OnStart(proxy IBufProxy)
	OnTriggle(proxy IBufProxy)
	OnEnd(proxy IBufProxy)
}

// IScriptMgr
// 提供脚本创建
type IScriptMgr interface {
	AddProvider(provider IScriptProvider)
	CreateBufScript(name string) IBufScript
}

// IScriptProvider
// 抽象golang,lua的脚本创建
type IScriptProvider interface {
	CreateBufScript(name string) IBufScript
}

// IBufScriptMgr
// for goprovider
type IBufScriptMgr interface {
	Create(name string) IBufScript
}

type ProviderCreator func(args ...any) IScriptProvider
