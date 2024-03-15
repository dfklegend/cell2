package sceneservice

// KickAllPlayers 将玩家都踢出去
func (s *Scene) KickAllPlayers() {
	uids := make([]int64, 0)
	for k, _ := range s.players {
		uids = append(uids, k)
	}

	for _, uid := range uids {
		s.KickPlayer(uid)
	}
}

func (s *Scene) KickPlayer(uid int64) {
	p := s.GetPlayer(uid)
	if p == nil {
		return
	}
	p.Kick()
}
