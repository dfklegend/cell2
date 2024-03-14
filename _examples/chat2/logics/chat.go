package logics

import (
	"fmt"
	"math/rand"
	"sort"

	"chat2/messages/clientmsg"
	"github.com/dfklegend/cell2/node/builtin/channel"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/common"
	"github.com/dfklegend/cell2/utils/logger"
)

const (
	MaxPlayerInRoom = 8
)

// ----

// ChatMgr
// 聊天具体的service
// 按登录 n人分一个房间
// 聊天房间内同步
// 输入/roll,掷骰子
// 离开后，其他人可以加入
type ChatMgr struct {
	// 不同房间
	idService *common.SerialIdService
	ns        *service.NodeService
	cs        *channel.Service

	rooms   map[uint32]*ChatRoom
	players map[string]*ChatPlayer
}

func NewChatService(
	ns *service.NodeService,
	cs *channel.Service) *ChatMgr {
	service := &ChatMgr{
		ns:        ns,
		cs:        cs,
		idService: common.NewSerialIdService(),
		rooms:     make(map[uint32]*ChatRoom),
		players:   make(map[string]*ChatPlayer),
	}
	return service
}

func (s *ChatMgr) Start() {
}

func (s *ChatMgr) findPlayer(uid string) *ChatPlayer {
	p, _ := s.players[uid]
	return p
}

func (s *ChatMgr) findRoom(roomId uint32) *ChatRoom {
	r, _ := s.rooms[roomId]
	return r
}

// TODO: 优先房间id比较小得
func (s *ChatMgr) findEmptyRoom() *ChatRoom {
	keys := make([]uint32, 0, len(s.rooms))

	for k, _ := range s.rooms {
		keys = append(keys, k)
	}

	//sort.Ints(keys)
	sort.SliceStable(keys, func(a, b int) bool {
		return keys[a] < keys[b]
	})

	for _, k := range keys {
		v := s.findRoom(k)
		if v.GetPlayerNum() < MaxPlayerInRoom {
			return v
		}
	}
	return nil
}

func (s *ChatMgr) GetPlayerRoom(uid string) *ChatRoom {
	p := s.findPlayer(uid)
	if p == nil {
		return nil
	}
	room := s.findRoom(p.RoomId)
	if room == nil {
		return nil
	}
	return room
}

func (s *ChatMgr) CreateRoom() *ChatRoom {
	id := s.idService.AllocId()
	s.rooms[id] = NewRoom(id, s.cs)
	return s.rooms[id]
}

func (s *ChatMgr) PlayerEnter(uid string, name string,
	frontId string, netId uint32) {
	p := NewPlayer()
	p.UId = uid
	p.Name = name
	p.FrontId = frontId
	p.NetId = netId

	s.players[uid] = p

	r := s.findEmptyRoom()
	if r == nil {
		r = s.CreateRoom()
	}

	// 推送房间members
	s.pushMembers(p, r)

	r.AddPlayer(p)
	// 广播进入房间
	newMsg := &clientmsg.OnNewUser{
		Name: p.GetDetailName(),
	}
	r.GetChannel().PushMessage("onNewUser", newMsg)

	// 显示系统信息
	s.showWelcome(p)
}

func (s *ChatMgr) pushMembers(p *ChatPlayer, r *ChatRoom) {
	cs := s.cs

	members := r.GetMembers()
	if len(members) == 0 {
		return
	}

	cs.PushMessageById(s.ns, p.FrontId, p.NetId, "onMembers",
		&clientmsg.OnMembers{
			Members: members,
		}, nil)
}

func (s *ChatMgr) showWelcome(p *ChatPlayer) {
	cs := s.cs

	str := fmt.Sprintf("welcome to %v room:%v", s.ns.Name, p.RoomId)
	logger.Log.Debugf(str)
	cs.PushMessageById(s.ns, p.FrontId, p.NetId, "onMessage",
		&clientmsg.ChatMsg{
			Name:    "system",
			Content: str,
		}, nil)
}

func (s *ChatMgr) PlayerLeave(uid string) {
	// find player
	p := s.findPlayer(uid)
	if p == nil {
		return
	}
	room := s.findRoom(p.RoomId)
	if room == nil {
		return
	}

	room.RemovePlayer(uid)
	leaveMsg := &clientmsg.OnUserLeave{
		Name: p.GetDetailName(),
	}
	room.GetChannel().PushMessage("onUserLeave", leaveMsg)

	delete(s.players, uid)
}

func (s *ChatMgr) OnChat(uid, content string) {
	p := s.findPlayer(uid)
	if p == nil {
		logger.Log.Errorf("can not find player: %v", uid)
		return
	}
	room := s.findRoom(p.RoomId)
	if room == nil {
		return
	}
	room.OnChat(p, content)
}

// ----

type ChatPlayer struct {
	UId     string
	Name    string
	FrontId string
	NetId   uint32
	RoomId  uint32
}

func NewPlayer() *ChatPlayer {
	return &ChatPlayer{}
}

func (p *ChatPlayer) GetDetailName() string {
	return fmt.Sprintf("%v %v.%v", p.Name, p.FrontId, p.NetId)
}

// ---------------

type ChatRoom struct {
	id      uint32
	players map[string]*ChatPlayer

	c  *channel.Channel
	cs *channel.Service
}

func NewRoom(id uint32, cs *channel.Service) *ChatRoom {
	return &ChatRoom{
		id:      id,
		players: make(map[string]*ChatPlayer),
		c:       cs.AddChannel(fmt.Sprintf("room%v", id)),
		cs:      cs,
	}
}

func (r *ChatRoom) Destroy() {
	r.cs.DeleteChannel(r.c.GetName())
}

func (r *ChatRoom) AddPlayer(p *ChatPlayer) {
	r.c.Add(p.FrontId, p.NetId)
	r.players[p.UId] = p
	p.RoomId = r.id
}

func (r *ChatRoom) RemovePlayer(uid string) {
	p, _ := r.players[uid]
	if p == nil {
		return
	}
	r.c.Leave(p.FrontId, p.NetId)
	delete(r.players, uid)
	p.RoomId = 0
}

func (r *ChatRoom) GetChannel() *channel.Channel {
	return r.c
}

func (r *ChatRoom) GetPlayerNum() int {
	return len(r.players)
}

func (r *ChatRoom) GetMembers() []string {
	members := make([]string, 0)
	for _, v := range r.players {
		members = append(members, v.GetDetailName())
	}
	return members
}

func (r *ChatRoom) OnChat(p *ChatPlayer, content string) {
	logger.Log.Debugf("OnChat:", content)
	if content == "/roll" {
		// do roll
		str := fmt.Sprintf("掷出了%v点(1-100)",
			rand.Intn(100))
		r.PushMessage(p, str)
		return
	}
	r.PushMessage(p, content)
}

func (r *ChatRoom) PushMessage(p *ChatPlayer, content string) {
	r.c.PushMessage("onMessage", &clientmsg.ChatMsg{
		Name:    p.UId,
		Content: content,
	})
}
