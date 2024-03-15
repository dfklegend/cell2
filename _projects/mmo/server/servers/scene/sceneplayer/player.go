package sceneplayer

import (
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/common"
	"github.com/dfklegend/cell2/utils/event/light"

	define2 "mmo/common/define"
	"mmo/common/entity"
	mymsg "mmo/messages"
	"mmo/messages/cproto"
	"mmo/modules/systems"
	"mmo/servers/scene/define"
	"mmo/utils"
)

// ScenePlayer
// 玩家角色
// 玩家数据所在
type ScenePlayer struct {
	ns *service.NodeService

	uid     int64
	frontId string
	netId   uint32
	logicId string

	online      bool
	timeOffline int64
	// 上一次发起logout的时间戳，避免重入
	lastTimeReqLogout int64
	logoutErrTimes    int
	// 致命错误，player最快速度移除自己
	fatalErr bool

	// avatar 和 camera
	avatar entity.EntityID
	camera entity.EntityID

	dirt  bool
	state define.ScenePlayerState

	scene       define.IScene
	sceneServer string
	sceneCfgId  int32
	sceneId     uint64

	systems *systems.Systems
	events  *light.EventCenter

	nextKeepAlive     int64
	keepAliveErrTimes int

	// 避免太频繁触发切换点切换
	lastExitTriggled int64
	// 正在切换
	sceneSwitching bool
}

func newScenePlayer(ns *service.NodeService, uid int64, frontId string, netId uint32, logicId string) *ScenePlayer {
	p := &ScenePlayer{
		ns:          ns,
		sceneServer: ns.Name,
		uid:         uid,
		frontId:     frontId,
		netId:       netId,
		logicId:     logicId,
		online:      true,
		state:       define.Init,
		events:      light.NewEventCenter(),
	}

	p.createSystems()
	return p
}

func (p *ScenePlayer) GetUId() int64 {
	return p.uid
}

func (p *ScenePlayer) GetNodeService() *service.NodeService {
	return p.ns
}

func (p *ScenePlayer) GetEvents() *light.EventCenter {
	return p.events
}

func (p *ScenePlayer) PushMsg(route string, msg interface{}) {
	app.PushMessageById(p.ns, p.frontId, p.netId, route, msg)
}

func (p *ScenePlayer) GetId() int64 {
	return p.uid
}

func (p *ScenePlayer) GetFrontId() string {
	return p.frontId
}

func (p *ScenePlayer) GetNetId() uint32 {
	return p.netId
}

func (p *ScenePlayer) GetLogicId() string {
	return p.logicId
}

func (p *ScenePlayer) UpdateScene(scene define.IScene, cfgId int32, sceneId uint64) {
	p.scene = scene
	p.sceneCfgId = cfgId
	p.sceneId = sceneId
}

func (p *ScenePlayer) LeaveCurScene() {
	if p.scene != nil {
		p.scene.PlayerLeave(p.uid)
		p.scene = nil
		p.sceneId = 0
	}
}

func (p *ScenePlayer) GetSceneId() uint64 {
	return p.sceneId
}

func (p *ScenePlayer) SetOnline(v bool) {
	p.online = v
}

func (p *ScenePlayer) IsOnline() bool {
	return p.online
}

func (p *ScenePlayer) SetDirt() {
	p.dirt = true
}

func (p *ScenePlayer) ClearDirt() {
	p.dirt = false
}

func (p *ScenePlayer) IsDirt() bool {
	return p.dirt
}

func (p *ScenePlayer) ChangeState(state define.ScenePlayerState) {
	p.state = state
}

func (p *ScenePlayer) GetState() define.ScenePlayerState {
	return p.state
}

func (p *ScenePlayer) IsState(state define.ScenePlayerState) bool {
	return p.state == state
}

func (p *ScenePlayer) SetAvatar(id entity.EntityID) {
	p.avatar = id
}

func (p *ScenePlayer) GetAvatar() entity.EntityID {
	return p.avatar
}

func (p *ScenePlayer) GetAvatarEntity() entity.IEntity {
	if p.avatar == entity.InvalidEntityId || p.scene == nil {
		return nil
	}
	return p.scene.GetEntity(p.avatar)
}

func (p *ScenePlayer) HasAvatar() bool {
	return p.avatar != entity.InvalidEntityId
}

func (p *ScenePlayer) ClearAvatar() {
	p.avatar = entity.InvalidEntityId
}

func (p *ScenePlayer) SetCamera(id entity.EntityID) {
	p.camera = id
}

func (p *ScenePlayer) GetCamera() entity.EntityID {
	return p.camera
}

func (p *ScenePlayer) HasCamera() bool {
	return p.camera != entity.InvalidEntityId
}

func (p *ScenePlayer) ClearCamera() {
	p.camera = entity.InvalidEntityId
}

func (p *ScenePlayer) Destroy() {
	p.systems.Destroy()
}

func (p *ScenePlayer) Update() {
	p.tryLogout()
	p.tryKeepAlive()
}

func (p *ScenePlayer) OnOffline() {
	// todo: 启动timer,开始定期检查，是否可以删除角色
	p.ns.GetLogger().Infof("%v offline", p.uid)
	p.SetOnline(false)
	p.timeOffline = common.NowMs()
	p.events.Publish("onOffline")

	// 需要移除绑定的camera
	if p.scene != nil {
		p.scene.DestroyCamera(p.uid)
	}
}

func (p *ScenePlayer) OnReOnline() {
	p.ns.GetLogger().Infof("%v reonline", p.uid)

	p.SetOnline(true)
	p.events.Publish("onReOnline")
}

func (p *ScenePlayer) OpenCamera() {
	// 打开场景内角色的摄像机
	p.notifyClientLoadScene(p.ns.Name, p.sceneCfgId, p.sceneId)
}

func (p *ScenePlayer) notifyClientLoadScene(sceneServer string, cfgId int32, sceneId uint64) {
	p.PushMsg("loadscene", &cproto.LoadScene{
		ServerId: sceneServer,
		CfgId:    cfgId,
		SceneId:  sceneId,
	})
}

func (p *ScenePlayer) tryKeepAlive() {
	if p.GetState() != define.Normal {
		return
	}
	now := common.NowMs()
	if now < p.nextKeepAlive {
		return
	}
	p.nextKeepAlive = now + define2.SceneKeepAliveMs

	ns := p.ns
	logger := ns.GetLogger()
	app.Request(ns, "logic.logicremote.keepalive", p.logicId, &mymsg.CheckPlayer2{
		UId:     p.uid,
		SceneId: p.sceneId,
	}, func(err error, ret interface{}) {
		ack := utils.TryGetNormalAck(ret)
		if err != nil || (ack != nil && ack.Code != int32(define2.Succ)) {
			logger.Errorf("keepalive failed! player: %v", p.uid)
			if err != nil {
				logger.Errorf("keepalive failed! player: %v, err: %v", p.uid, err)
			}
			p.onKeepAliveFailed()
			return
		}

		if define2.DebugSimulateSceneKeepAliveFailed {
			p.onKeepAliveFailed()
			return
		}
	})
}

func (p *ScenePlayer) onKeepAliveFailed() {
	logger := p.ns.GetLogger()
	p.keepAliveErrTimes++
	logger.Infof("%v onKeepAliveFailed times: %v", p.uid, p.keepAliveErrTimes)
	// 3次timeout, logic那边6次
	if p.keepAliveErrTimes > define2.SceneKeepAliveFailedTimes {
		p.SetFatalErr()
	}
}

func (p *ScenePlayer) Kick() {
	app.Kick(p.ns, p.frontId, p.netId, nil)
}
