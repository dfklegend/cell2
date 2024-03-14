package pomelo

import (
	"github.com/dfklegend/cell2/node/builtin/msgs"
	"github.com/dfklegend/cell2/node/client/impls"
	cs "github.com/dfklegend/cell2/node/client/session"
	"github.com/dfklegend/cell2/pomelonet/common/conn/message"
	pi "github.com/dfklegend/cell2/pomelonet/interfaces"
	"github.com/dfklegend/cell2/utils/sche"
)

// 	----------------
// 不同协议的实现，提供不同的impl
// 将客户端消息转化成ClientMsg同时，提供对应的serializer
// proxy of sessions
// pi.IClientSessionImpl -> sessions
type SessionsImpl struct {
	scheduler *sche.Sche
	sessions  *impls.ClientSessions
}

func NewSessionsImpl(scheduler *sche.Sche, sessions *impls.ClientSessions) *SessionsImpl {
	return &SessionsImpl{
		scheduler: scheduler,
		sessions:  sessions,
	}
}

func (s *SessionsImpl) translateSession(from pi.IClientSession) cs.IClientSession {
	// 目前是完全相等的，后续可以做proxy对象
	return from
}

func (s *SessionsImpl) OnSessionCreate(session pi.IClientSession) {
	s.scheduler.Post(func() {
		s.sessions.AddSession(s.translateSession(session))
	})
}

func (s *SessionsImpl) OnSessionClose(session pi.IClientSession) {
	s.scheduler.Post(func() {
		s.sessions.RemoveSession(s.translateSession(session))
	})
}

func (s *SessionsImpl) ProcessMessage(session pi.IClientSession, msg *message.Message) {
	// 处理消息
	// 转化成ClientMsg
	cmsg := &msgs.ClientMsg{}
	cmsg.SessionId = session.GetId()
	cmsg.ClientReqId = uint32(msg.ID)
	cmsg.Route = msg.Route
	cmsg.Data = msg.Data

	s.scheduler.Post(func() {
		s.sessions.ProcessMessage(s.translateSession(session), cmsg)
	})
}
