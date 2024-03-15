package units

import (
	"mmo/common/entity"
)

type MonsterUnitImpl struct {
	*BaseUnitImpl
}

func NewMonsterUnitImpl() *MonsterUnitImpl {
	return &MonsterUnitImpl{
		BaseUnitImpl: NewBaseUnitImpl(),
	}
}

func (p *MonsterUnitImpl) Init(entity entity.IEntity) {
	p.entity = entity
}

func (p *MonsterUnitImpl) Update() {
}

func (p *MonsterUnitImpl) OnDead() {
	e := p.entity
	e.GetWorld().DestroyEntity(e.GetId())
}
