package space

import (
	"mmo/common/entity"
	define2 "mmo/servers/scene/define"
	"mmo/servers/scene/space/factory"
)

func init() {
	factory.SetZoneCreator(func(args ...any) define2.ISpace {
		s := NewZoneSpace()
		s.Init(-define2.MaxWidth, -define2.MaxWidth, define2.MaxWidth, define2.MaxWidth,
			5)
		return s
	})
}

type ZoneSpace struct {
	entities map[entity.EntityID]*ZoneEntityInfo
	zones    []*Zone

	beginX   float32
	beginZ   float32
	endX     float32
	endZ     float32
	zoneSise float32

	zoneWidth  int
	zoneHeight int
}

func NewZoneSpace() *ZoneSpace {
	z := &ZoneSpace{
		entities: make(map[entity.EntityID]*ZoneEntityInfo),
	}
	return z
}

func (s *ZoneSpace) Init(beginX, beginZ, endX, endZ, zoneSize float32) {
	s.beginX = beginX
	s.beginZ = beginZ
	s.endX = endX
	s.endZ = endZ
	s.zoneSise = zoneSize

	s.zoneWidth = int((endX-beginX)/zoneSize) + 1
	s.zoneHeight = int((endZ-beginZ)/zoneSize) + 1

	s.zones = make([]*Zone, s.zoneWidth*s.zoneHeight)
	for i := 0; i < len(s.zones); i++ {
		s.zones[i] = NewZone()
	}
}

func (s *ZoneSpace) nToZoneN(n, begin, step float32, max int) int {
	zoneN := int((n - begin) / step)
	if zoneN < 0 {
		return 0
	}
	if zoneN >= max {
		return max - 1
	}
	return zoneN
}

func (s *ZoneSpace) xToZoneX(x float32) int {
	return s.nToZoneN(x, s.beginX, s.zoneSise, s.zoneWidth)
}

func (s *ZoneSpace) zToZoneZ(z float32) int {
	return s.nToZoneN(z, s.beginZ, s.zoneSise, s.zoneHeight)
}

func (s *ZoneSpace) xzToIndex(x, z float32) int {
	return s.zToZoneZ(z)*s.zoneWidth + s.xToZoneX(x)
}

func (s *ZoneSpace) AddEntity(id entity.EntityID, pos define2.Pos) {
	if s.entities[id] != nil {
		return
	}

	info := NewZoneInfo()
	info.Id = id
	info.Pos = pos

	s.entities[id] = info
	s.addToZone(info, pos.X, pos.Z)
}

func (s *ZoneSpace) addToZone(info *ZoneEntityInfo, x, z float32) {
	zoneX := s.xToZoneX(x)
	zoneZ := s.zToZoneZ(z)
	index := zoneZ*s.zoneWidth + zoneX

	zone := s.zones[index]
	zone.AddEntity(info)

	info.ZoneIndex = index
}

func (s *ZoneSpace) RemoveEntity(id entity.EntityID) {
	info := s.entities[id]
	if info == nil {
		return
	}

	index := info.ZoneIndex
	zone := s.zones[index]
	zone.RemoveEntity(info)

	delete(s.entities, id)
}

func (s *ZoneSpace) UpdateEntityPos(id entity.EntityID, pos define2.Pos) {
	info := s.entities[id]
	if info == nil {
		return
	}

	info.Pos = pos

	// check if index is same
	oldIndex := info.ZoneIndex
	newIndex := s.xzToIndex(pos.X, pos.Z)
	if oldIndex == newIndex {
		return
	}

	// remove from old
	if oldIndex < 0 {
		panic("unexpect")
	}
	zone := s.zones[oldIndex]
	succ := zone.RemoveEntity(info)
	if !succ {
		panic("unexpect")
	}

	s.zones[newIndex].AddEntity(info)
	info.ZoneIndex = newIndex
}

// SearchCircleTargets 圆形范围寻路
func (s *ZoneSpace) SearchCircleTargets(pos define2.Pos, radius float32, searcher define2.ISearcher) []entity.EntityID {
	// 根据范围查找zone上下限
	// 每个zone搜索一下

	xMinZone := s.xToZoneX(pos.X - radius)
	xMaxZone := s.xToZoneX(pos.X + radius)
	zMinZone := s.zToZoneZ(pos.Z - radius)
	zMaxZone := s.zToZoneZ(pos.Z + radius)
	for z := zMinZone; z <= zMaxZone; z++ {
		index := z*s.zoneWidth + xMinZone
		for x := xMinZone; x <= xMaxZone; x++ {
			s.zones[index].SearchCircleTargets(pos, radius, searcher)
			index++
		}
	}
	return searcher.MakeResults()
}
