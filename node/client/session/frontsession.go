package session

const (
	KeyUId      = "_ID"
	KeyNetId    = "_NetId"
	KeyServerId = "_ServerId"
)

// FrontSession
// 可以在session上设置一些数据
// 后端服务器可以将数据推送到前端设置，并且在后续获取
// 保留
//      _ID         Bind的唯一ID
//      _NetId      对应的frontSession id
//      _ServerId   对应的frontserver id

// 	Notice: 需要注意的是由于是json序列化过来，数字类型的数据可能被转成float64
type FrontSession struct {
	Session IClientSession
	Data    *SessionData
}

func NewFrontSession(serverId string, session IClientSession) *FrontSession {
	f := &FrontSession{
		Session: session,
		Data:    NewSessionData(),
	}

	f.Set(KeyServerId, serverId)
	f.Set(KeyNetId, session.GetId())
	return f
}

func (s *FrontSession) Reserve() {}
func (s *FrontSession) Handle()  {}

func (s *FrontSession) Bind(id string) {
	s.Data.Set(KeyUId, id)
}

func (s *FrontSession) GetID() string {
	return s.Get(KeyUId, "").(string)
}

func (s *FrontSession) GetNetId() uint32 {
	return s.Get(KeyNetId, 0).(uint32)
}

func (s *FrontSession) Set(k string, v interface{}) {
	s.Data.Set(k, v)
}

func (s *FrontSession) Get(k string, def interface{}) interface{} {
	return s.Data.Get(k, def)
}

func (s *FrontSession) PushSession(cb func(error)) {
	// doNothing
	if cb != nil {
		cb(nil)
	}
}

func (s *FrontSession) QuerySession(cb func(error)) {
	if cb != nil {
		cb(nil)
	}
}

func (s *FrontSession) IsSessionDataReady() bool {
	return true
}

func (s *FrontSession) ToJson() string {
	return s.Data.ToJsonStr()
}

func (s *FrontSession) FromJson(d string) {
	s.Data.FromJsonStr(d)
}

func (s *FrontSession) Kick() {
	s.Session.Close()
}

func (s *FrontSession) Lock() {
}

func (s *FrontSession) Unlock() {
}

func (s *FrontSession) IsClosed() bool {
	return s.Session.IsClosed()
}
