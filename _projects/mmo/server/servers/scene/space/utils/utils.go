package utils

import (
	"mmo/common/entity"
	"mmo/servers/scene/define"
	define1 "mmo/servers/scene/entity/define"
	"mmo/servers/scene/space/searchers"
)

func FindNearestEnemy(space define.ISpace, e entity.IEntity, radius float32) entity.EntityID {
	tran := e.GetComponent(define1.Transform).(define1.ITransform)
	if tran == nil {
		return 0
	}
	searcher := searchers.NewFindNearestEnemy(e)
	entities := space.SearchCircleTargets(tran.GetPos(), radius, searcher)
	if entities == nil {
		return 0
	}

	return entities[0]
}

// FindPlayers
// 返回范围内的Avatar对象
// ret: never be nil
func FindPlayers(space define.ISpace, e entity.IEntity, radius float32) []entity.EntityID {
	tran := e.GetComponent(define1.Transform).(define1.ITransform)
	if tran == nil {
		return nil
	}
	searcher := searchers.NewFindPlayers(e)
	return space.SearchCircleTargets(tran.GetPos(), radius, searcher)
}
