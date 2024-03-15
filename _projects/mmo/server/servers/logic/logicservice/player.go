package logicservice

import (
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/common"

	"mmo/common/define"
	mymsg "mmo/messages"
)

// LogicPlayer
// 维护角色跨服行为
// 没有数据需要存储
type LogicPlayer struct {
	// 所在gate
	FrontId string
	NetId   uint32

	uid  int64
	dirt bool

	// 战斗相关
	sceneCfgId     int
	sceneId        uint64
	sceneServer    string
	sceneEnterTime int64
	// 保证切线和切换场景互斥
	switchSceneLocked bool

	ns *service.NodeService

	// 需要定期收到scene的keepAlive消息
	// keepAlive失败，则认为玩家失效了，需要删除自己，并通知center
	lastKeepAliveChecked  int64
	keepAliveTimeoutTimes int
}

func NewPlayer(ns *service.NodeService) *LogicPlayer {
	p := &LogicPlayer{
		ns:                   ns,
		lastKeepAliveChecked: common.NowMs(),
	}
	return p
}

func (p *LogicPlayer) GetUId() int64 {
	return p.uid
}

func (p *LogicPlayer) Init(frontId string, netId uint32) {

	p.ns.GetLogger().Infof("logic player %v bind to: (%v, %v) -> (%v, %v)", p.uid,
		p.FrontId, p.NetId,
		frontId, netId)

	p.FrontId = frontId
	p.NetId = netId
}

func (p *LogicPlayer) Destroy() {
}

func (p *LogicPlayer) EnterWorld() {
}

func (p *LogicPlayer) IsLogouting() bool {
	return false
}

func (p *LogicPlayer) PushMsg(route string, msg interface{}) {
	app.PushMessageById(p.ns, p.FrontId, p.NetId, route, msg)
}

func (p *LogicPlayer) RefreshSceneKeepAlive() {
	if define.DebugSimulateSceneKeepAliveFailed {
		return
	}
	p.lastKeepAliveChecked = common.NowMs()
	p.keepAliveTimeoutTimes = 0
}

// 避免调试断电造成检查超时，改成多次检查
// 整体 3*2 = 6次SceneKeepAliveMs
// 晚于scene那边的timeout清理
func (p *LogicPlayer) checkKeepAlive() {
	now := common.NowMs()
	if now > p.lastKeepAliveChecked+2*int64(define.SceneKeepAliveMs) {
		p.keepAliveTimeoutTimes++
		p.lastKeepAliveChecked = now

		logger := p.ns.GetLogger()
		logger.Infof("checkKeepAlive: %v timeout times: %v", p.uid, p.keepAliveTimeoutTimes)
	}
}

func (p *LogicPlayer) IsKeepAliveTimeout() bool {
	return p.keepAliveTimeoutTimes >= define.LogicSceneKeepAliveTimeoutTimes
}

func (p *LogicPlayer) destroyForKeepAliveTimeout() {
	logger := p.ns.GetLogger()
	logger.Infof("destroyForKeepAliveTimeout: %v", p.uid)
	p.notifyCenterLogicAbnormalLogout()
	p.Destroy()
}

func (p *LogicPlayer) notifyCenterLogicAbnormalLogout() {
	ns := p.ns
	app.Request(ns, "center.centerremote.onabnormallogout", nil, &mymsg.OnLogout{
		UId: p.uid,
	}, func(err error, ack any) {

	})
}
