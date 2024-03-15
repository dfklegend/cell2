package clientmsg1

type EmptyArgReq struct {
}

type EmptyArg struct {
}

type NormalAck struct {
	Code int    `json:"code"`
	Err  string `json:"err"`
}

// 定义client消息
type QueryGateReq struct {
}

type QueryGateAck struct {
	Code int    `json:"code"`
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginAck struct {
	Code int32 `json:"code"`
	UId  int64 `json:"uid"`
}

type StartGame struct {
}

// ack NormalAck

// CharInfo 下发角色信息
type CharInfo struct {
	Name  string `json:"name"`
	Level int32  `json:"level"`
	Exp   int64  `json:"exp"`
	Money int64  `json:"money"`
}

type BattleLog struct {
	Log string `json:"log"`
}
