package client

/*
TODO
	考虑如何抽象来方便的支持其他协议
*/

import (
	"errors"

	"github.com/dfklegend/cell2/node/builtin/msgs"
	cs "github.com/dfklegend/cell2/node/client/session"
	"github.com/dfklegend/cell2/node/service"
)

type SessionId uint32

var (
	ErrorNoSession = errors.New("no session")
)

//	ISessionsHandler
// 	定义如何处理来自客户端的消息等
// 	设置到实际连接管理对象中去
type ISessionsHandler interface {
	Process(session *cs.FrontSession, msg *msgs.ClientMsg)
	OnSessionAdd(session *cs.FrontSession)
	OnSessionRemove(session *cs.FrontSession)
}

//	IClientSessions
// 	客户端session管理，提供接口用于管理
type IClientSessions interface {
	CloseSession(id SessionId)
	// PushMsg 推送消息
	PushMsg(msg *msgs.PushMsg)
	// Response 返回消息
	Response(response *msgs.Response)
}

type CBSessionOnClose func(*service.NodeService, *cs.FrontSession)
