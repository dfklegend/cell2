package common

import (
	"mmo/modules/lua"
)

type IScriptMgr interface {
}

// IWorld
// 战斗角色对外部的依赖
// 比如，获取其他对象，范围搜索等
type IWorld interface {
	GetTimeProvider() ITimeProvider
	GetWatcher() IWatcher
	GetDetailRecorder() IFightDetailRecorder

	GetLua() *lua.Service
	GetScriptMgr() IScriptMgr

	GetChar(id CharId) ICharacter
	// 其他，基于范围的搜索接口

}
