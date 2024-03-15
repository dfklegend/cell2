package components

import (
	"mmo/common/entity"
	"mmo/common/factory"
	common2 "mmo/modules/fight/common"
	"mmo/servers/scene/define"
)

// IUnitImpl 不同unit的行为差异
// 比如玩家和怪物
type IUnitImpl interface {
	Init(entity entity.IEntity)
	OnDestroy()
	Update()
	OnDead()
	OnKillTarget(tar common2.CharId, skillId common2.SkillId)
}

var UnitImplFactory = factory.NewIntFactory()

func CreateUnitImpl(unitType define.UnitType) IUnitImpl {
	ret := UnitImplFactory.Create(int(unitType))
	if ret == nil {
		return nil
	}
	return ret.(IUnitImpl)
}
