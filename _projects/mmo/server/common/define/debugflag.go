package define

// 定义一些调试开关
// 测试机制完善，缺省都是false

const (
	DebugForceSwitchLineSceneEnterTokenErr = false // 切线进入的token故意设置错误
	DebugForceSwitchLineSceneEnterFail1    = false // 强制切线enter时失败1
	DebugForceSwitchLineSceneEnterFail2    = false // 强制切线enter时失败2
	DebugSimulateSceneKeepAliveFailed      = false // 模拟场景keepalive失败
)
