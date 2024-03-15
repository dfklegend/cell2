package common

import (
	"github.com/dfklegend/cell2/utils/event/light"

	"mmo/modules/fight/attr"
)

// ICharacter
// id 为 entity id
// 所有可战斗角色的基础
// 		拥有属性
// 		可以释放技能
// 		身上可以驻留buf
// 		特殊属性
type ICharacter interface {
	IBufCtrl
	ISpecialStatusCtrl

	Init(id CharId, world IWorld, events *light.EventCenter)
	SetCharWatcher(watcher ICharWatcher)
	// Start 都设置完毕后
	Start()

	Destroy()

	GetId() CharId
	GetPos() Pos

	GetWorld() IWorld

	GetChar(id CharId) ICharacter
	GetEvents() *light.EventCenter
	Update()

	//	attr 属性相关

	AddAttrChangeWatcher(watcher IAttrChangeWatcher)
	RemoveAttrChangeWatcher(watcher IAttrChangeWatcher)

	VisitAttr(doFunc func(index int, item attr.IAttr))
	GetValue(index int) attr.Value

	GetIntValue(index int) int
	SetBaseValue(index int, v attr.Value)
	SetIntBaseValue(index int, v int)

	OffsetBase(index int, value attr.Value)
	OffsetIntBase(index int, value int)
	OffsetPercent(index int, off float32)

	GetLevel() int
	ChangeLevel(level int)
	GetSide() int
	SetHP(hp int)
	GetHP() int
	GetNormalAttackInterval() int32
	ApplyMPCost(mpcost int)

	ApplyDmg(dmg int) int
	//	attr end

	// -- IEquipable

	// AddEquipGroup equipgroup
	AddEquipGroup(group IEquipGroup)

	// equip slots

	SetEquip(index int, id EquipId)
	RemoveEquip(index int)

	Born()

	// -- IEquipable

	// -- skill

	GetTar() ICharacter
	GetTarId() CharId
	ClearTar()
	UpdateAttack(tar CharId)
	IsSkillRunning() bool

	GetSkillTable() ISkillTable
	// CastSkill 释放技能
	CastSkill(id SkillId, level int, src ICharacter, tar CharId)
	CallbackSkill(id SkillId, level int, src ICharacter, tar CharId)
	BreakSkill(src ICharacter, breakNormalAttack bool)
	//CastSkillAtPos(id SkillId, level int, tarPos Pos)
	// BackgroundSkill 后台技能
	//BackgroundSkill(id SkillId)

	PushSkillCD(id SkillId, prefireTime int32)
	OnSkillHit(dmg *DmgInstance)
	OnSkillBeHit(dmg *DmgInstance)
	OnSkillBroken(id SkillId)
	OnKillTarget(tar CharId, skillId SkillId)

	GetProxy() ICharProxy

	OnSpecialStatusChanged(id int, old bool, now bool)
	IsDead() bool
	IsInvincible() bool
}

type ICharProxy interface {
	GetId() int32
	GetTarId() CharId
	GetEvents() *light.EventCenter
	GoCallbackSkill(id SkillId, level int, src ICharProxy, tar CharId)
	AddBuf(id BufId, level int, stack int)
}

// ICharWatcher 可以收到char内部的一些反馈
type ICharWatcher interface {
	OnDead()
	OnKillTarget(tar CharId, skillId SkillId)
}

type IAttrChangeWatcher interface {
	OnAttrChanged(attrIndex int, oldV, newV attr.Value)
}
