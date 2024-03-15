package logic_old

import (
	"fmt"
	"time"

	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/event/light"
	l "github.com/dfklegend/cell2/utils/logger"

	mymsg "mmo/messages"
	"mmo/messages/cproto"
	"mmo/servers/logic_old/fightmode"
	"mmo/servers/logic_old/systems/bridge"
)

// Player player of logic_old
type Player struct {
	ns *service.NodeService

	// 所在gate
	FrontId string
	NetId   uint32

	UId int64

	// 简单信息
	Name  string
	Level int32
	Exp   int64
	Money int64

	dirt bool

	// 战斗相关
	sceneId        uint64
	sceneServer    string
	sceneEnterTime int64

	events *light.EventCenter

	logouting bool

	curFight fightmode.ISceneFightMode

	//systems *Systems
	// system interfaces
	charCard bridge.ICharCard
}

func NewPlayer(ns *service.NodeService) *Player {
	p := &Player{
		ns:     ns,
		events: light.NewEventCenter(),
	}
	//p.systems = newSystems(p)
	p.makeSystemInterfaces()
	return p
}

func (p *Player) GetUId() int64 {
	return p.UId
}

func (p *Player) GetEvents() *light.EventCenter {
	return p.events
}

func (p *Player) GetNodeService() *service.NodeService {
	return p.ns
}

func (p *Player) SetDirt() {
	p.dirt = true
}

func (p *Player) ClearDirt() {
	p.dirt = false
}

func (p *Player) IsDirt() bool {
	return p.dirt
}

func (p *Player) BeginLogout() {
	p.logouting = true
}

func (p *Player) IsLogouting() bool {
	return p.logouting
}

func (p *Player) Init(frontId string, netId uint32) {
	p.FrontId = frontId
	p.NetId = netId
}

func (p *Player) InitInfo() {
	p.initInfo()
}

func (p *Player) LoadInfo(info *mymsg.PlayerInfo) {
	// info can not be nil
	p.Name = info.Base.Name
	p.Level = info.Base.Level
	p.Exp = info.Base.Exp
	p.Money = info.Base.Money

	//p.systems.LoadInfo(info)
}

func (p *Player) initInfo() {
	l.L.Infof("player %v initinfo", p.UId)
	p.Level = 1

	//p.systems.InitInfo()
	p.SetDirt()
}

func (p *Player) MakeData() *mymsg.PlayerInfo {
	info := &mymsg.PlayerInfo{
		Base: &mymsg.BaseInfo{},
	}

	info.Base.UId = p.UId
	info.Base.Name = p.Name
	info.Base.Level = p.Level
	info.Base.Exp = p.Exp
	info.Base.Money = p.Money

	//p.systems.SaveInfo(info)
	return info
}

func (p *Player) Destroy() {
	//p.systems.Destroy()
}

func (p *Player) AddReward(money int, exp int) {
	p.Money += int64(money)
	p.Exp += int64(exp)
	p.SetDirt()
}

func (p *Player) EnterWorld() {
	l.L.Infof("player %v enter world", p.UId)

	// push infos
	p.PushMsg("enterbegin", &cproto.EmptyArg{})
	p.PushCharInfo()
	//p.systems.PushInfoToClient()
	p.PushMsg("enterend", &cproto.EmptyArg{})

	// enter world
	//p.systems.OnEnterWorld()
	p.events.Publish("onEnterWorld")
	// enter over

	//
	p.ns.GetRunService().GetTimerMgr().After(time.Second, func(args ...interface{}) {
		p.PushBattleLog(fmt.Sprintf("welcome to %v, running for %v s", p.ns.Name,
			app.Node.GetNodeCtrl().GetPassed()/1000))
	})
}

func (p *Player) PushCharInfo() {
	p.PushMsg("charinfo", &cproto.CharInfo{
		Name:  p.Name,
		Level: p.Level,
		Exp:   p.Exp,
		Money: p.Money,
		Cards: p.charCard.GetCards(),
	})
}

func (p *Player) PushMsg(route string, msg interface{}) {
	app.PushMessageById(p.ns, p.FrontId, p.NetId, route, msg)
}

func (p *Player) IsIdle() bool {
	return !p.IsInScene()
}
