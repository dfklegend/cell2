package units

import (
	"time"

	"github.com/dfklegend/cell2/utils/event/light"

	"mmo/common/entity"
	common2 "mmo/modules/fight/common"
	"mmo/servers/scene"
	"mmo/servers/scene/define"
	"mmo/servers/scene/entity/components"
	define2 "mmo/servers/scene/entity/define"
	"mmo/servers/scene/sceneplayer/systems/interfaces"
)

type PlayerUnitImpl struct {
	*BaseUnitImpl

	waitRelive bool
	scene      *sceneservice.Scene
	unit       *components.BaseUnit
	player     define.IPlayer
}

func NewPlayerUnitImpl() *PlayerUnitImpl {
	return &PlayerUnitImpl{
		BaseUnitImpl: NewBaseUnitImpl(),
	}
}

func (p *PlayerUnitImpl) Init(entity entity.IEntity) {
	p.BaseUnitImpl.Init(entity)

	p.scene = p.entity.GetWorld().GetContext().(*sceneservice.Scene)
	p.unit = p.entity.GetComponent(define2.BaseUnit).(*components.BaseUnit)
	p.player = p.entity.GetComponent(define2.Player).(*components.PlayerComponent).GetPlayer()
	p.bindEvents(true)
}

func (p *PlayerUnitImpl) OnDestroy() {
	p.BaseUnitImpl.OnDestroy()
	p.bindEvents(false)
}

func (p *PlayerUnitImpl) bindEvents(bind bool) {
	events := p.player.GetEvents()
	light.BindEventWithReceiver(bind, events, "onlevelup", p, p.onLevelUp)
}

func (p *PlayerUnitImpl) Update() {
}

func (p *PlayerUnitImpl) OnDead() {
	if p.waitRelive {
		return
	}
	p.waitRelive = true

	scene := p.scene
	timerMgr := scene.GetNodeService().GetRunService().GetTimerMgr()
	timerMgr.After(10*time.Second, func(args ...any) {
		p.onRelive()
	})
}

func (p *PlayerUnitImpl) onRelive() {
	p.waitRelive = false
	p.unit.Relive()
}

func (p *PlayerUnitImpl) OnKillTarget(tar common2.CharId, skillId common2.SkillId) {
	// 获取经验
	baseInfo := p.player.GetSystem(interfaces.BaseInfo).(interfaces.IBaseInfo)
	baseInfo.AddExp(10)
	baseInfo.PushInfoUpdate()
}

func (p *PlayerUnitImpl) onLevelUp(args ...any) {
	level := p.player.GetSystem(interfaces.BaseInfo).(interfaces.IBaseInfo).GetLevel()
	p.unit.GetChar().ChangeLevel(int(level))
}
