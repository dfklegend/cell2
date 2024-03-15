package center

const (
	// 	player status
	// logining和logouting两个状态不可打断
	// 提供长时间超时
	Init = iota
	Logining

	Logined // 登录成功，通过keepalive来规避异常状况(比如整个logic挂了)
	SwitchLine

	Logouting  // 登出过程中
	WaitRemove // 等待移除
	Abnormal   // 非正常状态，比如logined如果keepalive失败，可以转成此状态
)

/*
	loginning
		如果超时没有收到OnLogicOnline，则认为登录失败
	Logouting
		如果超时没有收到下线成功通知(超时时间较长)，为严重错误，可能导致玩家数据问题

*/

// 主要是为了避免重点事务的重入
// 上下线，切线等
const (
	TransactionInit       = iota
	TransactionLogin      // 登录事务，由客户端发起
	TransactionLogout     // 登出事务，由scene->logic->center发起
	TransactionReonline   // 断线重连，重新上线
	TransactionSwitchLine // 切线事务，由logic发起
)

const (
	TransactionLockTimeout      = 3 * 60 * 1000 // 常规事务超时
	TransactionLockLoginTimeout = 5 * 60 * 1000 // login超时
)
