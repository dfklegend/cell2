package chat

import (
	"fmt"
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/utils/common"
	mymsg "simplechat/messages"
)

const (
	MemberExpireTime int64 = 10 * 1000
)

type Member struct {
	ID  int32
	PID *actor.PID
	//Name         string
	TimeToExpire int64
}

// ----

type Room struct {
	ID      int32
	Token   int32
	Members []*Member

	roll    *RollCtrl
	service *Service
	players *PlayerMgr
}

func NewRoom() *Room {
	return &Room{
		Members: make([]*Member, 0),
		roll:    NewRoll(),
	}
}

func (r *Room) Init(s *Service, players *PlayerMgr) {
	r.service = s
	r.players = players
	r.roll.Init(s, r)
}

func (r *Room) GetPlayers() *PlayerMgr {
	return r.players
}

func (r *Room) AddMember(id int32, pid *actor.PID) {
	if r.FindMember(id) != nil {
		return
	}
	m := &Member{
		ID:  id,
		PID: pid,
		//Name:         name,
		TimeToExpire: common.NowMs() + MemberExpireTime,
	}
	r.Members = append(r.Members, m)
}

func (r *Room) FindMember(id int32) *Member {
	for _, v := range r.Members {
		if v.ID == id {
			return v
		}
	}
	return nil
}

func (r *Room) Broadcast(s *Service, msg interface{}) {
	for _, v := range r.Members {
		s.Request(v.PID, msg, nil)
	}
}

func (r *Room) Send(s *Service, id int32, msg interface{}) {
	m := r.FindMember(id)
	if m == nil {
		return
	}
	s.Request(m.PID, msg, nil)
}

func (r *Room) SendInfo(s *Service, id int32, str string) {
	r.Send(s, id, &mymsg.Chat{
		ID:   0,
		Name: "系统",
		Str:  str,
	})
}

func (r *Room) update() {
	r.roll.Update()
}

func (r *Room) removeExpired() bool {
	now := common.NowMs()
	dirt := false
	for i := 0; i < len(r.Members); {
		m := r.Members[i]
		if now < m.TimeToExpire {
			i++
			continue
		}

		dirt = true
		log.Printf("%v %v expired!", m.PID, m.ID)
		r.Members = append(r.Members[:i], r.Members[i+1:]...)
	}
	return dirt
}

func (r *Room) Report(s *Service) {
	// 向chatm汇报自己情况
	s.reportRoomStat(r.ID, int32(len(r.Members)))
}

func (r *Room) TryProcessCmd(id int32, str string) bool {
	if str == "/rollbegin" {
		r.roll.Start(id)
		return true
	}
	if str == "/roll" {
		r.roll.PlayerRoll(id)
		return true
	}
	if str == "/rollend" {
		r.roll.ReqEnd(id)
		return true
	}
	return false
}

func (r *Room) OnNewMemberJoin(id int32) {
	// 返回后，才能下发消息
	r.service.GetRunService().GetTimerMgr().After(time.Millisecond, func(args ...interface{}) {
		r.SendInfo(r.service, id, fmt.Sprintf("欢迎来到房间%v", r.ID))
		r.notifyMembers(id)
	})
}

func (r *Room) notifyMembers(id int32) {
	r.SendInfo(r.service, id, "Members:")
	for _, v := range r.Members {
		if v.ID == id {
			continue
		}
		r.SendInfo(r.service, id, fmt.Sprintf("  %v", r.players.GetPlayerName(v.ID)))
	}
}

// ----

type RoomMgr struct {
	service *Service
	players *PlayerMgr
	rooms   map[int32]*Room
}

func NewMgr() *RoomMgr {
	return &RoomMgr{
		players: NewPlayerMgr(),
		rooms:   make(map[int32]*Room),
	}
}

func (r *RoomMgr) Init(s *Service) {
	r.service = s
}

func (r *RoomMgr) Start() {
	r.service.GetRunService().GetTimerMgr().AddTimer(time.Second, func(args ...interface{}) {
		r.update()
	})
}

func (r *RoomMgr) findRoom(roomID int32) *Room {
	return r.rooms[roomID]
}

func (r *RoomMgr) CreateRoom(roomID int32, token int32) bool {
	room := r.findRoom(roomID)
	if room != nil {
		return false
	}
	room = NewRoom()
	room.ID = roomID
	room.Token = token
	r.rooms[roomID] = room

	room.Init(r.service, r.players)
	return true
}

func (r *RoomMgr) GetRoomNum() int {
	return len(r.rooms)
}

func (r *RoomMgr) GetPlayerNum() int {
	return r.players.GetPlayerNum()
}

func (r *RoomMgr) update() {
	// 检查member过期
	for _, v := range r.rooms {
		v.update()
		if v.removeExpired() {
			v.Report(r.service)
		}
	}
}

func (r *RoomMgr) Join(id int32, pid *actor.PID, name string, roomID int32, token int32) bool {
	room := r.findRoom(roomID)
	if room == nil {
		return false
	}
	if room.Token != token {
		return false
	}

	if room.FindMember(id) != nil {
		return false
	}

	room.AddMember(id, pid)
	r.players.AddPlayer(id, name, roomID)
	room.OnNewMemberJoin(id)
	return true
}

func (r *RoomMgr) Chat(id int32, str string) {
	p := r.players.GetPlayer(id)
	if p == nil {
		return
	}
	room := r.findRoom(p.RoomID)
	if room == nil {
		return
	}

	if room.TryProcessCmd(id, str) {
		return
	}

	room.Broadcast(r.service, &mymsg.Chat{
		ID:   id,
		Name: p.Name,
		Str:  str,
	})
}

func (r *RoomMgr) Ping(id int32) {
	p := r.players.GetPlayer(id)
	if p == nil {
		return
	}
	room := r.findRoom(p.RoomID)
	if room == nil {
		return
	}

	m := room.FindMember(id)
	if m == nil {
		return
	}
	m.TimeToExpire = common.NowMs() + MemberExpireTime
}

func (r *RoomMgr) ChangeName(id int32, name string) {
	p := r.players.GetPlayer(id)
	if p == nil {
		return
	}
	p.Name = name
}
