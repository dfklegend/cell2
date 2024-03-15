package searchers

import (
	"mmo/common/entity"
	"mmo/servers/scene/define"
	define2 "mmo/servers/scene/entity/define"
)

type FindPlayers struct {
	owner     entity.IEntity
	ownerId   entity.EntityID
	ownerUnit define2.IBaseUnit

	tars    []entity.EntityID
	curDist float32
}

func NewFindPlayers(owner entity.IEntity) define.ISearcher {
	f := &FindPlayers{
		owner:     owner,
		ownerId:   owner.GetId(),
		ownerUnit: owner.GetComponent(define2.BaseUnit).(define2.IBaseUnit),
		tars:      make([]entity.EntityID, 0),
	}
	return f
}

func (f *FindPlayers) Validate(id entity.EntityID, dist float32) bool {
	if id == f.ownerId {
		return false
	}

	// 合法性
	tar := f.owner.GetWorld().GetEntity(id)
	if tar == nil {
		return false
	}
	unit := tar.GetComponent(define2.BaseUnit).(define2.IBaseUnit)
	if unit == nil {
		return false
	}
	if unit.IsDead() {
		return false
	}

	// 玩家对象
	if unit.GetUnitType() != define.UnitAvatar {
		return false
	}

	return true
}

func (f *FindPlayers) AddCandidate(id entity.EntityID, dist float32) {
	f.tars = append(f.tars, id)
}

func (f *FindPlayers) MakeResults() []entity.EntityID {
	return f.tars
}
