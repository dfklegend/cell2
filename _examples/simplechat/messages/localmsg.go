package messages

// 定义一些本地消息
type ClientLogin struct {
	i int32
}

type ClientSay struct {
	Str string
}

type ClientNickname struct {
	Name string
}
