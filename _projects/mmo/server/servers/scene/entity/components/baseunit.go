package components

import (
	"math/rand"

	"github.com/dfklegend/cell2/utils/common"
	"github.com/dfklegend/cell2/utils/event/light"

	"mmo/common/entity"
	"mmo/messages/cproto"
	"mmo/modules/fight/builder"
	common2 "mmo/modules/fight/common"
	"mmo/servers/scene/define"
)

const AttackInterval int64 = 3000

type BaseUnit struct {
	*BaseSceneComponent

	unitType define.UnitType

	Name string

	character common2.ICharacter

	// 用于entity events
	// player的events主要用于system的交互
	events *light.EventCenter

	nextCanAttack int64

	impl IUnitImpl
}

func NewBaseUnit(ut define.UnitType) *BaseUnit {
	u := &BaseUnit{
		BaseSceneComponent: NewBaseSceneComponent(),
		unitType:           ut,
		events:             light.NewEventCenter(),
		nextCanAttack:      common.NowMs() + (rand.Int63() % 5000) + AttackInterval,
		impl:               CreateUnitImpl(ut),
	}

	u.character = builder.NewBuilder().Build()
	return u
}

func (u *BaseUnit) GetEvents() *light.EventCenter {
	return u.events
}

func (u *BaseUnit) GetChar() common2.ICharacter {
	return u.character
}

func (u *BaseUnit) OnPrepare() {
	u.BaseSceneComponent.OnPrepare()
	u.impl.Init(u.GetOwner())
	u.character.Init(u.GetOwner().GetId(), u.scene.GetWorldForFight(), u.events)
	u.character.SetCharWatcher(u)
	u.character.Start()
}

func (u *BaseUnit) OnStart() {
	u.bindEvents(true)
}

func (u *BaseUnit) Update() {
	u.impl.Update()
	u.character.Update()
}

func (u *BaseUnit) OnDestroy() {
	u.bindEvents(false)
	u.impl.OnDestroy()
	u.character.Destroy()
}

func (u *BaseUnit) bindEvents(bind bool) {
}

func (u *BaseUnit) GetUnitType() define.UnitType {
	return u.unitType
}

func (u *BaseUnit) CanAttack(now int64) bool {
	return now >= u.nextCanAttack
}

func (u *BaseUnit) ApplyDmg(dmg int32) {
	char := u.character
	if u.IsDead() {
		return
	}
	char.OffsetIntBase(common2.HP, -int(dmg))
	if char.GetHP() <= 0 {
		char.SetHP(0)
		u.OnDead()
	}
}

func (u *BaseUnit) OnDead() {
	u.impl.OnDead()
	u.events.Publish("onDead")
	u.GetScene().GetEvents().Publish("onUnitDead", u.GetOwner())
}

func (u *BaseUnit) IsDead() bool {
	return u.character.GetHP() == 0
}

func (u *BaseUnit) Relive() {
	u.character.SetHP(u.character.GetIntValue(common2.HPMax))
	// 广播复活消息
	u.events.Publish("onRelive")

	// 消息广播
	u.scene.PushViewMsg(define.Pos{}, "unitrelive", &cproto.UnitRelive{
		Id: int32(u.GetOwner().GetId()),
		HP: int32(u.character.GetHP()),
	})
}

func (u *BaseUnit) OnKillTarget(tar common2.CharId, skillId common2.SkillId) {
	// 添加经验
	if u.impl != nil {
		u.impl.OnKillTarget(tar, skillId)
	}
}

func (u *BaseUnit) IsSkillRunning() bool {
	return u.character.IsSkillRunning()
}

func (u *BaseUnit) StartSkill(skillId int32, tarId entity.EntityID) {
	u.character.CastSkill("英勇打击", 1, u.character, tarId)
}

func (u *BaseUnit) UpdateAttack(tarId entity.EntityID) {
	u.character.UpdateAttack(tarId)
}

func (u *BaseUnit) ClearTar() {
	u.character.ClearTar()
}
