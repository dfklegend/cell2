package space

import (
	"golang.org/x/exp/slices"

	"mmo/common/entity"
	define2 "mmo/servers/scene/define"
	"mmo/servers/scene/space/factory"
)

func init() {
	factory.SetNormalCreator(func(args ...any) define2.ISpace {
		return NewSimpleSpace()
	})
}

func Visit() {
}

// SimpleSpace
// 简单的space，只是挂起来
type SimpleSpace struct {
	entities map[entity.EntityID]*EntityInfo
	values   []*EntityInfo
}

func NewSimpleSpace() *SimpleSpace {
	return &SimpleSpace{
		entities: make(map[entity.EntityID]*EntityInfo),
	}
}

func (s *SimpleSpace) AddEntity(id entity.EntityID, pos define2.Pos) {
	if s.entities[id] != nil {
		s.entities[id].Pos = pos
		return
	}

	info := NewInfo()
	info.Id = id
	info.Pos = pos
	s.entities[id] = info

	s.values = append(s.values, info)
}

func (s *SimpleSpace) RemoveEntity(id entity.EntityID) {
	delete(s.entities, id)

	index := s.findIndex(id)
	if index != -1 {
		s.values = slices.Delete(s.values, index, index+1)
	}
}

func (s *SimpleSpace) findIndex(id entity.EntityID) int {
	for i := 0; i < len(s.values); i++ {
		info := s.values[i]
		if info.Id == id {
			return i
		}
	}
	return -1
}

// UpdateEntityPos 更新位置
func (s *SimpleSpace) UpdateEntityPos(id entity.EntityID, pos define2.Pos) {
	info := s.entities[id]
	if info == nil {
		return
	}
	info.Pos = pos
}

// SearchCircleTargets 圆形范围寻路
func (s *SimpleSpace) SearchCircleTargets(pos define2.Pos, radius float32, searcher define2.ISearcher) []entity.EntityID {
	//for k, v := range s.entities {
	//	dist := pos.Distance(v.Pos)
	//	if dist > radius {
	//		continue
	//	}
	//	if searcher.Validate(k, dist) {
	//		searcher.AddCandidate(k, dist)
	//	}
	//}
	//
	//return searcher.MakeResults()

	for i := 0; i < len(s.values); i++ {
		v := s.values[i]
		k := v.Id
		dist := pos.Distance(v.Pos)
		if dist > radius {
			continue
		}
		if searcher.Validate(k, dist) {
			searcher.AddCandidate(k, dist)
		}
	}

	return searcher.MakeResults()
}
