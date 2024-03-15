package logicservice

import (
	"time"

	"github.com/dfklegend/cell2/node/service"
)

type PlayerMgr struct {
	ns      *service.NodeService
	players map[int64]*LogicPlayer
}

func NewPlayerMgr(ns *service.NodeService) *PlayerMgr {
	return &PlayerMgr{
		ns:      ns,
		players: make(map[int64]*LogicPlayer),
	}
}

func (m *PlayerMgr) Start() {
	m.ns.GetRunService().GetTimerMgr().AddTimer(time.Second, m.onUpdate)
}

func (m *PlayerMgr) GetPlayer(uid int64) *LogicPlayer {
	return m.players[uid]
}

func (m *PlayerMgr) GetPlayerNum() int {
	return len(m.players)
}

func (m *PlayerMgr) CreatePlayer(uid int64) *LogicPlayer {
	p := NewPlayer(m.ns)
	p.uid = uid
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

func (m *PlayerMgr) Filter(cb func(p *LogicPlayer)) {
	for _, v := range m.players {
		cb(v)
	}
}

func (m *PlayerMgr) Kick(uid int64) {
}

func (m *PlayerMgr) onUpdate(args ...any) {
	m.checkKeepAliveTimeout()
}

func (m *PlayerMgr) checkKeepAliveTimeout() {

	var toRemove []*LogicPlayer
	for _, v := range m.players {
		if v.IsKeepAliveTimeout() {
			if toRemove == nil {
				toRemove = make([]*LogicPlayer, 0)
			}
			toRemove = append(toRemove, v)
		} else {
			v.checkKeepAlive()
		}
	}

	if toRemove == nil {
		return
	}
	for _, v := range toRemove {
		uid := v.uid
		v.destroyForKeepAliveTimeout()
		delete(m.players, uid)
	}
}
