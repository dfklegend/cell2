package chat

type Player struct {
	ID     int32
	RoomID int32
	Name   string
}

type PlayerMgr struct {
	players map[int32]*Player
}

func NewPlayerMgr() *PlayerMgr {
	return &PlayerMgr{
		players: make(map[int32]*Player),
	}
}

func (m *PlayerMgr) GetPlayer(id int32) *Player {
	return m.players[id]
}

func (m *PlayerMgr) GetPlayerName(id int32) string {
	p := m.GetPlayer(id)
	if p == nil {
		return "missing player"
	}
	return p.Name
}

func (m *PlayerMgr) GetPlayerNum() int {
	return len(m.players)
}

func (m *PlayerMgr) AddPlayer(id int32, name string, roomId int32) {
	m.players[id] = &Player{
		ID:     id,
		Name:   name,
		RoomID: roomId,
	}
}

func (m *PlayerMgr) RemovePlayer(id int32) {
	delete(m.players, id)
}
