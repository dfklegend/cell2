package searchers

import (
	"mmo/common/entity"
	common2 "mmo/modules/fight/common"
	"mmo/servers/scene/define"
	define1 "mmo/servers/scene/entity/define"
)

type FindNearestEnemy struct {
	owner     entity.IEntity
	ownerId   entity.EntityID
	ownerUnit define1.IBaseUnit

	enemy   entity.EntityID
	curDist float32
}

func NewFindNearestEnemy(owner entity.IEntity) define.ISearcher {
	f := &FindNearestEnemy{
		owner:     owner,
		ownerId:   owner.GetId(),
		ownerUnit: owner.GetComponent(define1.BaseUnit).(define1.IBaseUnit),
	}
	return f
}

func (f *FindNearestEnemy) Validate(id entity.EntityID, dist float32) bool {
	if id == f.ownerId {
		return false
	}

	// 合法性
	enemy := f.owner.GetWorld().GetEntity(id)
	if enemy == nil {
		return false
	}
	unit := enemy.GetComponent(define1.BaseUnit).(define1.IBaseUnit)
	if unit == nil {
		return false
	}
	if unit.IsDead() {
		return false
	}

	// 阵营
	if unit.GetChar().GetIntValue(common2.Side) == f.ownerUnit.GetChar().GetIntValue(common2.Side) {
		return false
	}

	if f.enemy == 0 {
		return true
	}
	return dist < f.curDist
}

func (f *FindNearestEnemy) AddCandidate(id entity.EntityID, dist float32) {
	f.enemy = id
	f.curDist = dist
}

func (f *FindNearestEnemy) MakeResults() []entity.EntityID {
	if f.enemy == 0 {
		return nil
	}
	return []entity.EntityID{f.enemy}
}
