package session

// frontsession, backsession，为backend服务器提供

// IServerSession
// 可以直接通过session，向客户端推送消息
// 向客户端推送消息，已知[serverId,netId]即可
type IServerSession interface {
	// Bind 绑定ID
	Bind(id string)
	GetID() string

	GetNetId() uint32

	// 	Set/Get设置值
	Set(k string, v interface{})
	Get(k string, def interface{}) interface{}
	// 	推送更新
	PushSession(cb func(error))
	ToJson() string
	FromJson(string)

	QuerySession(cb func(error))
	IsSessionDataReady() bool

	// Kick 踢下线
	Kick()

	// Lock
	// 互斥锁
	// 使用Bind,GetID,Set/Get
	// ToJson,FromJson
	Lock()
	Unlock()

	IsClosed() bool
}

// IClientSession
// 代表clientSession
type IClientSession interface {
	Reserve()

	GetId() uint32
	SetId(uint32)
	Push(route string, v interface{}) error
	// ResponseMID 返回
	ResponseMID(mid uint, v interface{}, e error) error
	Close()

	IsClosed() bool
}
