package logic_old

import (
	"github.com/dfklegend/cell2/node/service"
)

// PlayerMgr
type PlayerMgr struct {
	ns      *service.NodeService
	players map[int64]*Player
}

func NewPlayerMgr(ns *service.NodeService) *PlayerMgr {
	return &PlayerMgr{
		ns:      ns,
		players: make(map[int64]*Player),
	}
}

func (m *PlayerMgr) Start() {
}

func (m *PlayerMgr) GetPlayer(uid int64) *Player {
	return m.players[uid]
}

func (m *PlayerMgr) GetPlayerNum() int {
	return len(m.players)
}

func (m *PlayerMgr) CreatePlayer(uid int64) *Player {
	p := NewPlayer(m.ns)
	p.UId = uid
	m.players[uid] = p
	return p
}

func (m *PlayerMgr) DestroyPlayer(uid int64) {
	p := m.GetPlayer(uid)
	if p == nil {
		return
	}
	p.Destroy()
	delete(m.players, uid)
}

func (m *PlayerMgr) Filter(cb func(p *Player)) {
	for _, v := range m.players {
		cb(v)
	}
}

func (m *PlayerMgr) Kick(uid int64) {
	ReqOffline(m.ns, m, uid, func(err bool) {
		//
	})
}
