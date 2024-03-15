package sceneplayer

import (
	"time"

	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/event/light"
	"github.com/dfklegend/cell2/utils/logger/interfaces"

	"mmo/servers/scene/define"
)

type PlayerMgr struct {
	ns      *service.NodeService
	players map[int64]*ScenePlayer
	logger  interfaces.Logger
}

func NewPlayerMgr(ns *service.NodeService) *PlayerMgr {
	return &PlayerMgr{
		ns:      ns,
		players: make(map[int64]*ScenePlayer),
	}
}

func (m *PlayerMgr) Start() {
	m.logger = m.ns.GetLogger()
	m.ns.GetRunService().GetTimerMgr().AddTimer(time.Millisecond, m.update)
	events := m.ns.GetOwner().(define.ISceneService).GetEvents()
	light.BindEventWithReceiver(true, events, "onFatalErr", m, m.onFatalErr)
}

func (m *PlayerMgr) update(args ...any) {
	m.updatePlayers()
}

func (m *PlayerMgr) GetPlayer(uid int64) *ScenePlayer {
	return m.players[uid]
}

func (m *PlayerMgr) GetPlayerNum() int {
	return len(m.players)
}

func (m *PlayerMgr) CreatePlayer(uid int64, frontId string, netId uint32, logicId string) *ScenePlayer {
	player := newScenePlayer(m.ns, uid, frontId, netId, logicId)
	m.players[uid] = player
	return player
}

func (m *PlayerMgr) destroyPlayer(p *ScenePlayer) {
	m.logger.Infof("destroyPlayer: %v", p.uid)
	p.Destroy()
	delete(m.players, p.GetId())
}

func (m *PlayerMgr) updatePlayers() {
	var players []*ScenePlayer
	for _, v := range m.players {
		v.Update()
		if v.IsState(define.WaitMgrDelete) {
			if players == nil {
				players = make([]*ScenePlayer, 0)
				players = append(players, v)
			}
		}
	}

	if players == nil {
		return
	}
	for _, v1 := range players {
		m.destroyPlayer(v1)
	}
}

func (m *PlayerMgr) Filter(cb func(p *ScenePlayer)) {
	for _, v := range m.players {
		cb(v)
	}
}

func (m *PlayerMgr) Kick(uid int64) {
}

func (m *PlayerMgr) onFatalErr(args ...any) {
	logger := m.logger
	logger.Infof("player mgr onFatalErr")
	m.Filter(func(p *ScenePlayer) {
		p.SetFatalErr()
	})
}
