package units

import (
	"mmo/common/entity"
	common2 "mmo/modules/fight/common"
)

type BaseUnitImpl struct {
	entity entity.IEntity
}

func NewBaseUnitImpl() *BaseUnitImpl {
	return &BaseUnitImpl{}
}

func (p *BaseUnitImpl) Init(entity entity.IEntity) {
	p.entity = entity
}

func (p *BaseUnitImpl) OnDestroy() {

}

func (p *BaseUnitImpl) Update() {
}

func (p *BaseUnitImpl) OnDead() {

}

func (p *BaseUnitImpl) OnKillTarget(tar common2.CharId, skillId common2.SkillId) {
}
