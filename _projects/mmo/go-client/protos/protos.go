package protos

type EmptyArgReq struct {
}

type NormalAck struct {
	Code   int    `json:"code"`
	Result string `json:"result"`
}

type QueryGateReq struct {
}

type QueryGateAck struct {
	Code int    `json:"code"`
	IP   string `json:"ip"`
	Port string `json:"port"`
}

type Hello struct {
	Msg    string `json:"msg"`
	Number int    `json:"Number"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginAck struct {
	Code int32 `json:"code"`
	UId  int64 `json:"uid"`
}
