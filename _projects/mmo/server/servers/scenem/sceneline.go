package scenem

import (
	"math/rand"
	"sort"

	"golang.org/x/exp/slices"
)

// ISceneLine 同一个配置多条线
type ISceneLine interface {
	GetCfgId() int32
	GetSceneId() uint64
	GetLineId() int32
}

type SceneLine struct {
	cfgId   int32
	sceneId uint64
	lineId  int32
}

func NewSceneLine(cfgId int32, sceneId uint64, lineId int32) *SceneLine {
	return &SceneLine{
		cfgId:   cfgId,
		sceneId: sceneId,
		lineId:  lineId,
	}
}

func (s *SceneLine) GetCfgId() int32 {
	return s.cfgId
}

func (s *SceneLine) GetSceneId() uint64 {
	return s.sceneId
}

func (s *SceneLine) GetLineId() int32 {
	return s.lineId
}

type SceneLines struct {
	// sorted
	lines []*SceneLine
}

func NewLines() *SceneLines {
	return &SceneLines{
		lines: make([]*SceneLine, 0),
	}
}

// FineIdleLineId
// 返回一个最小可用的lineId
func (l *SceneLines) FineIdleLineId() int32 {
	var idleId int32

	idleId = 0
	for _, v := range l.lines {
		if v.lineId != idleId {
			break
		}
		idleId++
	}
	return idleId
}

func (l *SceneLines) add(line *SceneLine) {
	l.lines = append(l.lines, line)

	lines := l.lines
	sort.Slice(lines, func(a, b int) bool {
		return lines[a].lineId < lines[b].lineId
	})
}

func (l *SceneLines) findIndex(lineId int32) int {
	return slices.IndexFunc(l.lines, func(line *SceneLine) bool {
		return lineId == line.lineId
	})
}

func (l *SceneLines) remove(lineId int32) {
	index := l.findIndex(lineId)
	if index == -1 {
		return
	}
	l.lines = slices.Delete(l.lines, index, index+1)
}

func (l *SceneLines) GetNum() int {
	return len(l.lines)
}

func (l *SceneLines) RandGetScene() uint64 {
	// 随机取一个
	if l.GetNum() == 0 {
		return 0
	}
	index := rand.Int() % l.GetNum()
	return l.lines[index].sceneId
}
