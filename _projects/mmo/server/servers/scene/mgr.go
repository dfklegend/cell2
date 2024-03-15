package sceneservice

import (
	"time"

	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/common"
	l "github.com/dfklegend/cell2/utils/logger"

	"mmo/common/define"
	mymsg "mmo/messages"
	"mmo/modules/csv"
	"mmo/modules/fight/script"
	"mmo/modules/fight/script/builder"
	"mmo/modules/fightscripts/bridge"
	"mmo/modules/lua"
	define2 "mmo/servers/scene/define"
)

// SceneMgr 场景管理
type SceneMgr struct {
	scenes    map[uint64]*Scene
	ns        *service.NodeService
	lua       *lua.Service
	scriptMgr script.IScriptMgr

	nextRefresh     int64
	refreshErrTimes int
}

func NewSceneMgr() *SceneMgr {
	return &SceneMgr{
		scenes: make(map[uint64]*Scene),
	}
}

func (m *SceneMgr) SetLua(lua *lua.Service) {
	m.lua = lua
}

func (m *SceneMgr) GetSceneNum() int {
	return len(m.scenes)
}

func (m *SceneMgr) Start(ns *service.NodeService) {
	m.ns = ns
	m.initScriptMgr()

	ns.GetRunService().GetTimerMgr().AddTimer(10*time.Millisecond, m.update)
}

func (m *SceneMgr) initScriptMgr() {
	mgr := builder.CreateScriptMgr()
	m.scriptMgr = mgr

	lua := m.lua

	mgr.AddProvider(script.CreateGoProvider(bridge.GetBufMgr()))
	mgr.AddProvider(script.CreateLuaProvider(lua.GetEnvData()))
}

func (m *SceneMgr) update(args ...any) {
	m.tryRefreshToSceneM()

	var toRemove []uint64
	for k, v := range m.scenes {
		if v.IsOver() {
			if toRemove == nil {
				toRemove = make([]uint64, 0)
			}
			toRemove = append(toRemove, k)
		} else {
			v.Update()
		}
	}

	if toRemove != nil {
		for _, v := range toRemove {
			m.freeScene(v)
		}
	}
}

func (m *SceneMgr) tryRefreshToSceneM() {
	now := common.NowMs()
	if now < m.nextRefresh {
		return
	}
	m.nextRefresh = now + define.SceneToSceneMKeepAlive

	app.Request(m.ns, "scenem.remote.refresh", nil, &mymsg.SMRefresh{
		ServiceId: m.ns.Name,
		SceneNum:  int32(m.GetSceneNum()),
	}, func(err error, ret any) {
		// 如果连接失败，超过一定次数，移除所有场景
		m.onRefreshRet(err == nil)
	})
}

func (m *SceneMgr) onRefreshRet(succ bool) {
	if succ {
		m.refreshErrTimes = 0
		return
	}
	m.refreshErrTimes++
	if m.refreshErrTimes >= 5 {
		m.onFatalErr()
		m.refreshErrTimes = 0
	}
}

func (m *SceneMgr) onFatalErr() {
	m.forceCloseAllScenes()
	service := m.ns.GetOwner().(define2.ISceneService)
	service.GetEvents().Publish("onFatalErr")
}

func (m *SceneMgr) getLogicFromCfgId(cfgId int32) string {
	// 后续根据配置
	entry := csv.Scene.GetEntry(int(cfgId))
	if entry == nil {
		return "standard"
	}
	return "standard"
}

func (m *SceneMgr) AllocScene(cfgId int32, sceneId uint64, token int) bool {
	if m.scenes[sceneId] != nil {
		l.L.Errorf("already has scene: %v", sceneId)
		return false
	}
	scene := NewScene(m.lua, m.scriptMgr)
	m.scenes[sceneId] = scene

	scene.Init(m.ns, sceneId, m.getLogicFromCfgId(cfgId), token, cfgId)
	scene.Start()
	return true
}

func (m *SceneMgr) freeScene(sceneId uint64) {
	scene := m.scenes[sceneId]
	if scene == nil {
		l.L.Errorf("freeScene error find no scene: %v", sceneId)
		return
	}

	m.ns.GetLogger().Infof("free scene: %v", sceneId)
	scene.Destroy()
	delete(m.scenes, sceneId)
}

func (m *SceneMgr) FindScene(sceneId uint64) *Scene {
	return m.scenes[sceneId]
}

func (m *SceneMgr) forceCloseAllScenes() {
	for _, v := range m.scenes {
		v.SetFatalErr()
	}
}
