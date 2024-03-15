package units

import (
	"mmo/common/entity"
)

type StaticUnitImpl struct {
	*BaseUnitImpl
}

func NewStaticUnitImpl() *StaticUnitImpl {
	return &StaticUnitImpl{
		BaseUnitImpl: NewBaseUnitImpl(),
	}
}

func (p *StaticUnitImpl) Init(entity entity.IEntity) {
	p.entity = entity
}

func (p *StaticUnitImpl) Update() {
}

func (p *StaticUnitImpl) OnDead() {
	e := p.entity
	e.GetWorld().DestroyEntity(e.GetId())
}
