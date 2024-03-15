package sceneservice

import (
	l "github.com/dfklegend/cell2/utils/logger"

	"mmo/servers/scene/define"
)

// PlayerEnter 玩家进入
func (s *Scene) PlayerEnter(player define.IPlayer) bool {
	logger := s.ns.GetLogger()
	if s.players[player.GetId()] != nil {
		logger.Errorf("already has player: %v in scene: %v", player.GetId(), s.sceneId)
		return false
	}
	l.L.Infof("PlayerEnter: %v", player.GetId())

	s.players[player.GetId()] = player
	s.logic.PlayerEnter(player)
	player.UpdateScene(s, s.cfgId, s.sceneId)
	return true
}
