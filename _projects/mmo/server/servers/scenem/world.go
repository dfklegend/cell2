package scenem

import (
	"github.com/dfklegend/cell2/utils/common"
	l "github.com/dfklegend/cell2/utils/logger"
)

// IWorld 管理场景
type IWorld interface {
	GetSceneListByCfgId(cfgId int32) []ISceneLine
	NewSceneLine(cfgId int32) ISceneLine
	FreeSceneLine(cfgId int32, lineId int32)
}

type World struct {
	mgr *SceneServiceMgr
	// 某个cfgId的场景有几条线(副本)
	sceneLines map[int32]*SceneLines
	// sceneId: 每个具体的场景
	scenes       map[uint64]*SceneObj
	publicScenes *PublicScenes
}

func NewWorld() *World {
	return &World{
		sceneLines:   make(map[int32]*SceneLines),
		scenes:       make(map[uint64]*SceneObj),
		publicScenes: NewPublicScenes(),
	}
}

func (w *World) Init(mgr *SceneServiceMgr) {
	w.mgr = mgr
	w.publicScenes.Init(mgr)
}

func (w *World) Start() {
}

func (w *World) SpawnPublicScene() {
}

func (w *World) NewSceneLine(cfgId int32, sceneId uint64) *SceneLine {
	lines := w.sceneLines[cfgId]
	if lines == nil {
		lines = NewLines()
		w.sceneLines[cfgId] = lines
	}

	lineId := lines.FineIdleLineId()
	line := NewSceneLine(cfgId, sceneId, lineId)
	lines.add(line)
	return line
}

func (w *World) FreeSceneLine(cfgId int32, lineId int32) {
	lines := w.sceneLines[cfgId]
	if lines == nil {
		return
	}
	lines.remove(lineId)
}

func (w *World) GetLineNum(cfgId int32) int {
	lines := w.sceneLines[cfgId]
	if lines == nil {
		return 0
	}
	return lines.GetNum()
}

func (w *World) OnSceneCreateSucc(obj *SceneObj) {
	line := w.NewSceneLine(obj.CfgId, obj.SceneId)
	obj.LineId = line.lineId
	obj.TimeStart = common.NowMs()

	w.scenes[obj.SceneId] = obj
}

func (w *World) OnSceneEnd(sceneId uint64) {
	logger := w.mgr.ns.GetLogger()
	logger.Infof("OnSceneEnd: %v", sceneId)
	o := w.scenes[sceneId]
	if o == nil {
		l.L.Warnf("OnSceneEnd miss scene: %v", sceneId)
		return
	}

	w.FreeSceneLine(o.CfgId, o.LineId)
	delete(w.scenes, sceneId)
}

func (w *World) TrySpawnPublicScenes() {
	w.publicScenes.Start()
}

func (w *World) QueryScenes(maxNum int) []uint64 {
	scenes := make([]uint64, 0)
	for k, _ := range w.scenes {
		if len(scenes) >= maxNum {
			break
		}

		scenes = append(scenes, k)
	}

	return scenes
}

func (w *World) GetScene(sceneId uint64) *SceneObj {
	return w.scenes[sceneId]
}

// OnServiceLost
// 如果某个Service丢失
func (w *World) OnServiceLost(serviceId string) {
	// 移除所有该线对应的scene
	var toRemove []uint64
	for k, v := range w.scenes {
		if v.ServiceId != serviceId {
			continue
		}
		if toRemove == nil {
			toRemove = make([]uint64, 0)
		}
		toRemove = append(toRemove, k)
	}

	if toRemove == nil {
		return
	}

	for _, v := range toRemove {
		w.OnSceneEnd(v)
	}
}

func (w *World) ReqSceneByCfgId(cfgId int32) *SceneObj {
	lines := w.sceneLines[cfgId]
	if lines == nil {
		return nil
	}

	sceneId := lines.RandGetScene()
	if sceneId == 0 {
		return nil
	}
	return w.GetScene(sceneId)
}
