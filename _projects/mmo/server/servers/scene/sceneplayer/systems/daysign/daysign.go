package daysign

import (
	"reflect"

	"mmo/common/define"
	"mmo/common/factory"
	mymsg "mmo/messages"
	"mmo/messages/cproto"
	"mmo/modules/systems"
	sysfactory "mmo/servers/scene/sceneplayer/systems/factory"
	"mmo/servers/scene/sceneplayer/systems/interfaces"
	"mmo/utils"
)

func init() {
	name := interfaces.DaySign
	sysfactory.GetFactory().RegisterFunc(name, func(args ...any) factory.IObject {
		return newSystem(name)
	})
}

func Visit() {}

// 每日签到
// 每天可以请求签到一次
type system struct {
	*systems.BaseSystem

	lastDaySigned int32 // 从1970而来日子的index
}

func newSystem(name string) systems.ISystem {
	return &system{
		BaseSystem: systems.NewBaseSystem(name),
	}
}

func (s *system) OnCreate() {
	s.bindEvents(true)

	s.RegisterCmdHandler("sign", reflect.TypeOf(&cproto.TestAdd{}), s.onHandleSign)

	logger := s.GetPlayer().GetNodeService().GetLogger()
	logger.Infof("DaySign.OnCreate")

}

func (s *system) OnDestroy() {
}

func (s *system) bindEvents(bind bool) {
	//events := s.player.GetEvents()
}

func (s *system) InitData() {
}

func (s *system) LoadData(info *mymsg.PlayerInfo) {
	i := info.DaySign

	if i == nil {
		return
	}
	s.lastDaySigned = i.LastDaySigned
}

func (s *system) SaveData(info *mymsg.PlayerInfo) {
	info.DaySign = &mymsg.DaySign{}

	to := info.DaySign
	from := s
	to.LastDaySigned = from.lastDaySigned
}

func (s *system) OnEnterWorld(switchLine bool) {
}

func (s *system) PushInfoToClient() {
	s.PushCmd("init", &cproto.DaySign{
		Signed: s.lastDaySigned == utils.GetAbsoluteDay(),
	})
}

func (s *system) onHandleSign(srcArgs any, cb func(ret any, code int32)) {
	// 判断今天是否已经完成签到
	nowDay := utils.GetAbsoluteDay()
	if nowDay == s.lastDaySigned {
		cb(nil, int32(define.AlreadyDone))
		return
	}

	s.lastDaySigned = nowDay
	// 给奖励
	baseInfo := s.GetPlayer().GetSystem(interfaces.BaseInfo).(interfaces.IBaseInfo)
	baseInfo.AddMoney(100)
	baseInfo.PushSystemInfo(0, "签到成功")
	baseInfo.PushInfoUpdate()
	cb(nil, int32(define.Succ))
}
