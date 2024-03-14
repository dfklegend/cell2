package session

import (
	"encoding/json"
	"errors"

	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/builtin/msgs"
	"github.com/dfklegend/cell2/node/service"
	l "github.com/dfklegend/cell2/utils/logger"
)

// BackSession
// 后端服务处理请求时，可以设置数据到FrontSession，也可以获取之前设置的数据
// 需要知道session所在serverId和NetId
// 优化: 出于性能的考虑，SessionData 需要query才能获取
// 绑定的Id必然会传递过来
// 大部分接口不需要SessionData
type BackSession struct {
	ns       *service.NodeService
	ServerId string
	NetId    uint32
	dirt     bool
	Data     *SessionData
	NewData  *SessionData

	hasSessionData bool
}

//
// 	也可以直接创建一个backSession，用于push session
func NewBackSession(ns *service.NodeService, serverId string, netId uint32, ID string) *BackSession {
	v := NewBS()
	InitBackSession(v, ns, serverId, netId, ID)
	return v
}

func NewBS() *BackSession {
	return &BackSession{
		ns:   nil,
		Data: NewSessionData(),
		// new data saved, will push
		NewData: NewSessionData(),
	}
}

func InitBackSession(s *BackSession, ns *service.NodeService, serverId string, netId uint32, ID string) {
	s.ns = ns
	s.ServerId = serverId
	s.NetId = netId

	s.Data.Set(KeyUId, ID)
}

func ResetBackSession(s *BackSession) {
	s.ns = nil
	s.Data.Reset()
	s.NewData.Reset()
}

func CloneBackSession(from *BackSession) *BackSession {
	bs := NewBS()
	InitBackSession(bs, from.ns, from.ServerId, from.NetId, from.GetID())
	return bs
}

func (s *BackSession) Reserve() {}
func (s *BackSession) Handle()  {}

func (s *BackSession) Bind(id string) {
	s.Set(KeyUId, id)
}

func (s *BackSession) GetID() string {
	return s.Get(KeyUId, "").(string)
}

func (s *BackSession) GetNetId() uint32 {
	return s.NetId
}

func (s *BackSession) Set(k string, v interface{}) {
	s.NewData.Set(k, v)
	s.setDirt()
}

func (s *BackSession) setDirt() {
	s.dirt = true
}

func (s *BackSession) Get(k string, def interface{}) interface{} {
	if s.NewData.Has(k) {
		return s.NewData.Get(k, def)
	}
	return s.Data.Get(k, def)
}

// PushSession 避免推送了过多的数据
func (s *BackSession) PushSession(cb func(error)) {
	// 推送给目标
	if !s.dirt {
		if cb != nil {
			cb(nil)
		}
		return
	}
	s.dirt = false

	pid := app.GetServicePID(s.ServerId)
	if pid == nil {
		l.Log.Errorf("can not find service: %v", s.ServerId)
		if cb != nil {
			cb(errors.New("no service"))
		}
		return
	}

	msg := &msgs.PushSession{
		SessionId:   s.NetId,
		SessionData: s.NewData.ToJson(),
	}
	s.ns.RequestEx(pid, "sys.pushsession", msg, func(err error, ret any) {
		if cb != nil {
			cb(err)
		}
	})
}

func (s *BackSession) QuerySession(cb func(error)) {
	app.QuerySession(s.ns, s.ServerId, s.NetId, func(err error, raw any) {
		if err != nil {
			l.Log.Errorf("Query session error: %v", err)
			if cb != nil {
				cb(err)
			}
			return
		}

		msg, _ := raw.(*msgs.QuerySessionAck)
		if msg == nil {
			return
		}

		s.FromJson(string(msg.SessionData))
		if cb != nil {
			cb(nil)
		}
	})
}

func (s *BackSession) IsSessionDataReady() bool {
	return s.hasSessionData
}

func (s *BackSession) ToJson() string {
	// 合成一个全的
	m1 := s.Data.GetMap()
	m2 := s.NewData.GetMap()

	newMap := make(map[string]interface{})
	mergeMap(newMap, m1)
	mergeMap(newMap, m2)

	v, _ := json.Marshal(newMap)
	return string(v)
}

func mergeMap(m1, m2 map[string]interface{}) {
	for k, v := range m2 {
		m1[k] = v
	}
}

func (s *BackSession) FromJson(d string) {
	s.Data.FromJsonStr(d)
	s.ServerId = s.Get(KeyServerId, "n").(string)
	if s.Data.Has(KeyNetId) {
		// json转interface value时，number被转成float64
		s.NetId = uint32(s.Get(KeyNetId, 0).(float64))
	}

	s.dirt = false
}

func (s *BackSession) Kick() {
	app.Request(s.ns, "x.sys.kick", s.ServerId, &msgs.Kick{
		SessionId: s.NetId,
	}, nil)
}

func (s *BackSession) Lock() {
}

func (s *BackSession) Unlock() {
}

func (s *BackSession) IsClosed() bool {
	return false
}
