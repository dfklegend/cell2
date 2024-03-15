package space

import (
	"golang.org/x/exp/slices"

	"mmo/common/entity"
	define2 "mmo/servers/scene/define"
)

type ZoneEntityInfo struct {
	Id        entity.EntityID
	Pos       define2.Pos
	ZoneIndex int
}

func NewZoneInfo() *ZoneEntityInfo {
	return &ZoneEntityInfo{
		ZoneIndex: -1,
	}
}

// Zone 区域
type Zone struct {
	values []*ZoneEntityInfo
}

func NewZone() *Zone {
	return &Zone{
		values: make([]*ZoneEntityInfo, 0),
	}
}

func (z *Zone) AddEntity(info *ZoneEntityInfo) {
	z.values = append(z.values, info)
}

func (z *Zone) RemoveEntity(info *ZoneEntityInfo) bool {
	index := z.findIndex(info)
	if index == -1 {
		return false
	}
	z.values = slices.Delete(z.values, index, index+1)
	return true
}

func (z *Zone) findIndex(info *ZoneEntityInfo) int {
	return slices.IndexFunc(z.values, func(v *ZoneEntityInfo) bool {
		return v == info
	})
}

func (z *Zone) SearchCircleTargets(pos define2.Pos, radius float32, searcher define2.ISearcher) {
	for i := 0; i < len(z.values); i++ {
		v := z.values[i]
		k := v.Id
		dist := pos.Distance(v.Pos)
		if dist > radius {
			continue
		}
		if searcher.Validate(k, dist) {
			searcher.AddCandidate(k, dist)
		}
	}
}
