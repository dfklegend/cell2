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
	Log    int    `json:"Log"`
}
