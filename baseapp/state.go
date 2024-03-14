package baseapp

//	AppState
const (
	State0        int = iota // 初始
	StatePrepared            // 调用过了prepard
	StateStarting            // 启动中
	StateNormal              // 启动正常
	StateStoping             // 关闭中
	StateStopped             // 成功关闭
)
