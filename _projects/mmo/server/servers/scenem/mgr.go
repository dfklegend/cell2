package scenem

import (
	"math/rand"
	"time"

	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/common"

	"mmo/common/define"
)

type SceneServiceMgr struct {
	ns          *service.NodeService
	nextSceneId uint64

	// 场景服务状态
	services map[string]*SceneServiceStat

	world *World
}

func NewMgr(ns *service.NodeService) *SceneServiceMgr {
	return &SceneServiceMgr{
		ns:          ns,
		nextSceneId: 1,
		services:    make(map[string]*SceneServiceStat),
		world:       NewWorld(),
	}
}

func (m *SceneServiceMgr) GetNodeService() *service.NodeService {
	return m.ns
}

func (m *SceneServiceMgr) Start() {
	m.world.Init(m)
	m.world.Start()
	m.ns.GetRunService().GetTimerMgr().AddTimer(time.Second, m.onUpdate)
}

func (m *SceneServiceMgr) AllocScene(cfgId int32) *SceneObj {
	scene := &SceneObj{}

	idleService := m.FindIdleService()
	if idleService == "" {
		return nil
	}

	scene.CfgId = cfgId
	scene.SceneId = m.allocSceneId()
	scene.ServiceId = idleService
	scene.Token = rand.Intn(10000)

	return scene
}

func (m *SceneServiceMgr) FindIdleService() string {
	var idlest string
	var curWeight float32
	for k, v := range m.services {
		if !v.Working {
			continue
		}
		weight := v.GetBusyWeight()
		if idlest == "" || weight < curWeight {
			idlest = k
			curWeight = weight
		}
	}
	return idlest
}

func (m *SceneServiceMgr) allocSceneId() uint64 {
	v := m.nextSceneId
	m.nextSceneId++
	return v
}

func (m *SceneServiceMgr) OnServiceRefresh(serviceId string, sceneNum int) {
	service := m.services[serviceId]
	if service == nil {
		service = &SceneServiceStat{
			Working: true,
		}
		m.services[serviceId] = service
	} else {
		service.Working = true
	}

	service.ActiveSceneNum = sceneNum
	service.LastActiveTime = common.NowMs()
	service.ActiveFailedTimes = 0
	m.CheckToSpawnPublicScenes()
}

func (m *SceneServiceMgr) OnSceneCreateSucc(obj *SceneObj) {
	m.ns.GetLogger().Infof("scene %v-%v create succ", obj.ServiceId, obj.SceneId)
	m.world.OnSceneCreateSucc(obj)
}

func (m *SceneServiceMgr) OnSceneEnd(sceneId uint64) {
	m.ns.GetLogger().Infof("scene %v end", sceneId)
	m.world.OnSceneEnd(sceneId)
}

// CheckToSpawnPublicScenes 有场景服务连接后，等30s,再刷
func (m *SceneServiceMgr) CheckToSpawnPublicScenes() {
	m.world.TrySpawnPublicScenes()
}

func (m *SceneServiceMgr) QueryScenes(maxNum int) []uint64 {
	return m.world.QueryScenes(maxNum)
}

func (m *SceneServiceMgr) GetScene(sceneId uint64) *SceneObj {
	return m.world.GetScene(sceneId)
}

func (m *SceneServiceMgr) onUpdate(args ...any) {
	m.updateWorkingState()
}

func (m *SceneServiceMgr) updateWorkingState() {
	now := common.NowMs()
	for k, v := range m.services {
		if v.Working && now >= v.LastActiveTime+3*define.SceneToSceneMKeepAlive {
			m.onServiceKeepAliveFailed(k, v)
		}
	}
}

func (m *SceneServiceMgr) onServiceKeepAliveFailed(serviceId string, stat *SceneServiceStat) {
	logger := m.ns.GetLogger()
	stat.ActiveFailedTimes++
	stat.LastActiveTime = common.NowMs()

	logger.Infof("onServiceKeepAliveFailed: %v times: %v", serviceId, stat.ActiveFailedTimes)
	if stat.ActiveFailedTimes > 3 {
		m.onServiceLost(serviceId, stat)
	}
}

func (m *SceneServiceMgr) onServiceLost(serviceId string, stat *SceneServiceStat) {
	logger := m.ns.GetLogger()
	logger.Infof("scene service: %v lost", serviceId)
	stat.Working = false
	m.world.OnServiceLost(serviceId)
}

// ReqSceneByCfgId
// 请求可进入的场景
func (m *SceneServiceMgr) ReqSceneByCfgId(cfgId int32) *SceneObj {
	logger := m.ns.GetLogger()
	obj := m.world.ReqSceneByCfgId(cfgId)

	if obj == nil {
		logger.Infof("can not find scene, cfgId: %v", cfgId)
	}
	return obj
}
