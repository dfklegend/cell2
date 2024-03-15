package baseinfo

import (
	"mmo/common/factory"
	mymsg "mmo/messages"
	"mmo/messages/cproto"
	"mmo/modules/systems"
	sysfactory "mmo/servers/scene/sceneplayer/systems/factory"
	"mmo/servers/scene/sceneplayer/systems/interfaces"
)

func init() {
	name := interfaces.BaseInfo
	sysfactory.GetFactory().RegisterFunc(name, func(args ...any) factory.IObject {
		return newSystem(name)
	})
}

func Visit() {}

// 基础系统
// 存储一些基础数据
type system struct {
	*systems.BaseSystem
	interfaces.IBaseInfo

	// 玩家数据
	// 简单信息
	name       string
	level      int32
	exp        int64
	money      int64
	loginTimes int32
}

func newSystem(name string) systems.ISystem {
	return &system{
		BaseSystem: systems.NewBaseSystem(name),
	}
}

func (s *system) OnCreate() {
	s.bindEvents(true)

	logger := s.GetPlayer().GetNodeService().GetLogger()
	logger.Infof("IBaseInfo.OnCreate")
}

func (s *system) OnDestroy() {
}

func (s *system) bindEvents(bind bool) {
	//events := s.player.GetEvents()
}

func (s *system) PushInfoToClient() {
	s.GetPlayer().PushMsg("charinfo", &cproto.CharInfo{
		Name:  s.name,
		Level: s.level,
		Exp:   s.exp,
		Money: s.money,
	})
}

func (s *system) InitData() {
	s.name = ""
	s.level = 1
	s.exp = 0
	s.money = 0
	s.loginTimes = 0
}

func (s *system) LoadData(info *mymsg.PlayerInfo) {
	b := info.Base

	s.name = b.Name
	s.level = b.Level
	s.exp = b.Exp
	s.money = b.Money
	s.loginTimes = b.LoginTimes
}

func (s *system) SaveData(info *mymsg.PlayerInfo) {
	info.Base = &mymsg.BaseInfo{}

	b := info.Base
	b.UId = s.GetPlayer().GetUId()
	b.Name = s.name
	b.Level = s.level
	b.Exp = s.exp
	b.Money = s.money
	b.LoginTimes = s.loginTimes
}

func (s *system) OnEnterWorld(switchLine bool) {
	if !switchLine {
		s.loginTimes++
		s.money++

		s.GetPlayer().SetDirt()
	}
}

func (s *system) GetLoginTimes() int32 {
	return s.loginTimes
}

func (s *system) AddExp(exp int) {
	// 增加经验，升级
	if exp < 0 {
		return
	}
	s.exp += int64(exp)
	s.GetPlayer().SetDirt()

	if s.exp >= 100 {
		s.exp -= 100
		s.level++
		s.onLevelup()
	}
	s.GetPlayer().GetNodeService().GetLogger().Infof("%v addExp %v-%v, level: %v",
		s.GetPlayer().GetUId(), exp, s.exp, s.level)
}

func (s *system) GetLevel() int32 {
	return s.level
}

func (s *system) onLevelup() {
	// 抛事件
	s.GetPlayer().GetEvents().Publish("onlevelup", s.level)
	s.GetPlayer().GetNodeService().GetLogger().Infof("%v levelup : %v", s.GetPlayer().GetUId(), s.level)
}

func (s *system) AddMoney(v int) {
	if v < 0 {
		return
	}
	s.money += int64(v)
	s.GetPlayer().SetDirt()
}

func (s *system) SubMoney(v int) bool {
	if v < 0 {
		return false
	}

	if s.money < int64(v) {
		return false
	}

	s.money -= int64(v)
	s.GetPlayer().SetDirt()
	return true
}

func (s *system) PushSystemInfo(t int32, info string) {
	s.PushCmd("info", &cproto.SystemInfo{
		Info: info,
		Type: t,
	})
}

func (s *system) PushInfoUpdate() {
	// 简化推送
	s.PushCmd("updatecharinfo", &cproto.CharInfo{
		Name:  s.name,
		Level: s.level,
		Exp:   s.exp,
		Money: s.money,
	})
}
