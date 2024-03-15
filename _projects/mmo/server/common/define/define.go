package define

import (
	"github.com/dfklegend/cell2/actorex/service"
)

// ---- code ----

// ErrorCode 定义常规返回码
type ErrorCode int32

// 和Service Code对应起来
// 1000以下保留

const (
	Succ ErrorCode = iota
	AlreadyDone
)

const (
	ErrBegin ErrorCode = iota + service.CodeUserBegin
	// ErrSystemBusy 系统忙，请稍后再试
	ErrSystemBusy
	ErrReserved ErrorCode = service.CodeUserBegin + 99
)

// login code (1100-xx)
const (
	ErrAlreadyOnline         ErrorCode = iota + service.CodeUserBegin + 100
	ErrAuthFailed                      // 认证失败
	ErrLogicCannotFindPlayer           // logic上找不到玩家
	ErrWithStr                         // Err, 错误由str说明
	ErrFaild                           // 失败，不需要更多的说明
)

// ---- code ----

// scene keep alive
// 控制时间, logic晚于scene处理
const (
	SceneKeepAliveMs                = 15 * 1000 // scene向logic keepAlive间隔
	LogicSceneKeepAliveTimeoutTimes = 6
	SceneKeepAliveFailedTimes       = 3
)

const (
	SceneToSceneMKeepAlive = 1000 // scene服务向sceneM keepAlive
)
