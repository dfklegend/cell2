package impls

import (
	"github.com/dfklegend/cell2/node/builtin/msgs"
	"github.com/dfklegend/cell2/node/client"
	cs "github.com/dfklegend/cell2/node/client/session"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/common"
	"github.com/dfklegend/cell2/utils/logger"
)

// 	IClientSessions
// 	逻辑在拥有者(nodeservice)routine执行，由SessionsImpl post到目标routine
// 	SessionData
type ClientSessions struct {
	ns        *service.NodeService
	serverId  string
	sessions  map[client.SessionId]*cs.FrontSession
	idService *common.SerialIdService
	handler   client.ISessionsHandler

	onCloseCB client.CBSessionOnClose
	// 用户可以定制一个踢人，比如先发个消息告知
	kickHandler IKickHandler

	nextCanLogSessionMissLog int64
	skippedLogCount          int
}

func NewClientSessions(name string) *ClientSessions {
	return &ClientSessions{
		serverId:  name,
		sessions:  make(map[client.SessionId]*cs.FrontSession),
		idService: common.NewSerialIdService(),
	}
}

func (s *ClientSessions) SetNodeService(ns *service.NodeService) {
	s.ns = ns
}

func (s *ClientSessions) SetHandler(handler client.ISessionsHandler) {
	s.handler = handler
}

func (s *ClientSessions) SetKickHandler(handler IKickHandler) {
	s.kickHandler = handler
}

func (s *ClientSessions) SetOnCloseHandler(handler client.CBSessionOnClose) {
	s.onCloseCB = handler
}

func (s *ClientSessions) VisitSession(visitor func(session *cs.FrontSession)) {
	if visitor == nil {
		return
	}

	for _, v := range s.sessions {
		visitor(v)
	}
}

func (s *ClientSessions) allocId() client.SessionId {
	return client.SessionId(s.idService.AllocId())
}

func (s *ClientSessions) AddSession(session cs.IClientSession) {
	id := s.allocId()

	session.SetId(uint32(id))
	fs := cs.NewFrontSession(s.serverId, session)

	s.sessions[id] = fs

	if s.handler != nil {
		s.handler.OnSessionAdd(fs)
	}

}

func (s *ClientSessions) RemoveSession(session cs.IClientSession) {

	id := client.SessionId(session.GetId())
	fs := s.sessions[id]
	if fs == nil {
		logger.Log.Errorf("remove a session not exsit: %v", id)
		return
	}

	delete(s.sessions, id)

	if s.handler != nil {
		s.handler.OnSessionRemove(fs)
	}

	if s.onCloseCB != nil {
		s.onCloseCB(s.ns, fs)
	}

}

func (s *ClientSessions) findSession(netId uint32) *cs.FrontSession {

	id := client.SessionId(netId)
	return s.sessions[id]
}

func (s *ClientSessions) ProcessMessage(session cs.IClientSession, msg *msgs.ClientMsg) {
	//logger.Log.Infof("msg: %v", msg.Route)

	if s.handler != nil {
		s.handler.Process(s.findSession(session.GetId()), msg)
	}
}

//	PushMsg 下发消息
func (s *ClientSessions) PushMsg(msg *msgs.PushMsg) {

	for _, v := range msg.Ids {
		session := s.sessions[client.SessionId(v)]
		if session == nil {
			s.onSessionMissed(v, msg)
			continue
		}

		session.Session.Push(msg.Route, msg.Data)
	}
}

// 避免刷太多
func (s *ClientSessions) onSessionMissed(sessionId uint32, msg *msgs.PushMsg) {
	// 1秒只输出一次
	now := common.NowMs()
	if now < s.nextCanLogSessionMissLog {
		s.skippedLogCount++
		return
	}

	logger.L.Errorf("skipped session missed log: %v", s.skippedLogCount)
	logger.L.Errorf("can not find session: %v %v", sessionId, msg.Route)

	s.nextCanLogSessionMissLog = now + 1000
	s.skippedLogCount = 0
}

//	in owner routine
func (s *ClientSessions) PushSession(id uint32, data []byte) {
	session := s.findSession(id)
	if session == nil {
		return
	}
	session.Data.UpdateFromJson(data)
}

func (s *ClientSessions) Kick(id uint32) {
	session := s.findSession(id)
	if session == nil {
		return
	}
	logger.Log.Infof("req kick :%v", id)

	if s.kickHandler == nil {
		// 缺省踢人
		session.Kick()
	} else {
		s.kickHandler.HandleKick(s.ns, s, id)
	}
}

// DoKick 实际踢人
func (s *ClientSessions) DoKick(id uint32) {
	session := s.findSession(id)
	if session == nil {
		return
	}
	logger.Log.Infof("ClientSessions do kick :%v", id)
	session.Kick()
}

func (s *ClientSessions) GetSession(id uint32) *cs.FrontSession {
	return s.findSession(id)
}

// 	----------------
// SessionsComponent
type SessionsComponent struct {
	*service.BaseComponent
	sessions *ClientSessions
}

func NewSessionsComponent(sessions *ClientSessions) *SessionsComponent {
	return &SessionsComponent{
		BaseComponent: service.NewBaseComponent(),
		sessions:      sessions,
	}
}

func (s *SessionsComponent) GetSessions() *ClientSessions {
	return s.sessions
}

func (s *SessionsComponent) OnAdd() {
	s.sessions.SetNodeService(s.GetNodeService())
}
