package sceneservice

// PlayerLeave 玩家离开场景
func (s *Scene) PlayerLeave(uid int64) {
	logger := s.ns.GetLogger()
	logger.Infof("scene :%v PlayerLeave: %v", s.sceneId, uid)

	player := s.players[uid]
	if player == nil {
		return
	}

	s.logic.PlayerLeave(uid)
	delete(s.players, uid)
}
