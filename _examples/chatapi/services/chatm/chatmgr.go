package chatmgr

import (
	"errors"
	"log"

	mymsg "chatapi/messages"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/utils/common"
	"github.com/dfklegend/cell2/utils/runservice"
)

/*
	管理房间模块
	要求房间服务器创建 房间
*/

const (
	maxPlayerInRoom = 3
)

type Room struct {
	ChatID string
	RoomID int32
	Token  int32

	// 人数
	// 考虑，如何保证和房间服务实际的人数相等
	PlayerNum int32
	MaxNum    int32

	// 超过时间，房间移除
	Timeout int64
}

// 	每隔一段时间汇报状态
type ChatService struct {
	ChatID        string
	RoomNum       int32
	PlayerNum     int32
	LastRefreshed int64
}

type Mgr struct {
	rs         *runservice.StandardRunService
	nextRoomID int32

	rooms    map[int32]*Room
	services map[string]*ChatService
}

func NewMgr() *Mgr {
	return &Mgr{
		nextRoomID: 1,
		rooms:      make(map[int32]*Room),
		services:   make(map[string]*ChatService),
	}
}

func (m *Mgr) Init(service *runservice.StandardRunService) {
	m.rs = service
}

func (m *Mgr) allocRoomID() int32 {
	m.nextRoomID++
	return m.nextRoomID
}

func (m *Mgr) OnServiceInfo(chatID string, roomNum, playerNum int32) {
	service := m.services[chatID]
	if service == nil {
		service = &ChatService{
			ChatID: chatID,
		}
		m.services[chatID] = service
	}

	service.RoomNum = roomNum
	service.PlayerNum = playerNum
	service.LastRefreshed = common.NowMs()
}

func (m *Mgr) selectService() *ChatService {
	if len(m.services) == 0 {
		return nil
	}

	// select min load
	var chatId string
	var minLoad int32
	minLoad = -1
	for _, v := range m.services {
		if minLoad == -1 || minLoad > v.RoomNum {
			minLoad = v.RoomNum
			chatId = v.ChatID
		}
	}
	return m.services[chatId]
}

//	登录
func (m *Mgr) Login(service *Service, msg *mymsg.CMReqLogin, cbFunc func(error, any)) {
	// 1. 找找有没有空房间
	// 2. 没有空房间，创建一个
	// 3. 有空房间，房间信息

	room := m.findEmptyRoom()
	if room == nil {
		// create room
		m.loginCreateRoom(service, msg, cbFunc)
		return
	}

	room.PlayerNum++
	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.CMAckLogin{
		ChatServiceId: room.ChatID,
		RoomID:        room.RoomID,
		Token:         room.Token,
	})
}

func (m *Mgr) loginCreateRoom(service *Service, msg *mymsg.CMReqLogin, cbFunc func(error, any)) {
	cs := m.selectService()
	if cs == nil {
		apientry.CheckInvokeCBFunc(cbFunc, errors.New("no chat service"), nil)
		return
	}

	// 要求对方创建房间
	pid := app.GetServicePID(cs.ChatID)
	if pid == nil {
		apientry.CheckInvokeCBFunc(cbFunc, errors.New("service pid is nil"), nil)
		return
	}

	var roomID int32
	var token int32

	roomID = m.allocRoomID()
	token = 1234

	service.RequestEx(pid, "chat.createroom", &mymsg.CSCreateRoom{
		RoomID: roomID,
		Token:  token,
	}, func(err error, raw interface{}) {
		if err != nil {
			log.Printf("create room failed! err:%v \n", err)
			apientry.CheckInvokeCBFunc(cbFunc, err, nil)
			return
		}

		room := &Room{
			ChatID: cs.ChatID,
			RoomID: roomID,
			Token:  token,
			MaxNum: maxPlayerInRoom,
		}
		m.rooms[roomID] = room

		apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.CMAckLogin{
			ChatServiceId: cs.ChatID,
			RoomID:        roomID,
			Token:         token,
		})
	})
}

func (m *Mgr) findEmptyRoom() *Room {
	for _, v := range m.rooms {
		if v.PlayerNum+1 < v.MaxNum {
			return v
		}
	}
	return nil
}

func (m *Mgr) updateRoomStat(roomId int32, playerNum int32) {
	r := m.rooms[roomId]
	if r == nil {
		return
	}
	r.PlayerNum = playerNum

	log.Printf("%v num-> %v \n", roomId, playerNum)
}
