package sceneservice

import (
	"github.com/dfklegend/astar"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/event/light"

	"mmo/common/entity"
	"mmo/common/fightutils"
	"mmo/messages"
	"mmo/modules/csv"
	"mmo/modules/csv/entry"
	common2 "mmo/modules/fight/common"
	"mmo/modules/fight/script"
	"mmo/modules/lua"
	"mmo/servers/scene/blockgraph"
	"mmo/servers/scene/define"
	logic2 "mmo/servers/scene/logic"
)

// Scene 一个场景
//
type Scene struct {
	ns *service.NodeService

	sceneId  uint64
	token    int
	cfgId    int32
	initData define.ISceneInitData

	logic logic2.ISceneLogic
	space define.ISpace

	// 玩家对象
	// 来自logic
	players map[int64]define.IPlayer
	cameras map[entity.EntityID]define.ICamera

	// entities
	world  entity.IWorld
	events *light.EventCenter

	timeProvider  common2.ITimeProvider
	recorder      common2.IFightDetailRecorder
	worldForFight common2.IWorld
	lua           *lua.Service

	startTime int64
	destroyed bool
	fatalErr  bool

	width  int
	height int

	entry *entry.Scene
	graph *blockgraph.Graph

	finder      *astar.PathFinder
	pathContext *PathContext
}

func NewScene(lua *lua.Service, mgr script.IScriptMgr) *Scene {
	s := &Scene{
		world:        entity.NewWorld(),
		players:      map[int64]define.IPlayer{},
		cameras:      map[entity.EntityID]define.ICamera{},
		events:       light.NewEventCenter(),
		timeProvider: fightutils.NewTimeProvider(),
		lua:          lua,
	}

	//s.recorder = detailrecorder.NewRecorder(s.timeProvider)
	s.worldForFight = fightutils.NewWorldForFight(s.world, s.timeProvider,
		newFightWatcher(s), s.recorder, s.lua, mgr)
	return s
}

func (s *Scene) Init(ns *service.NodeService, sceneId uint64, logicType string, token int, cfgId int32) {
	s.ns = ns
	s.sceneId = sceneId
	s.token = token
	s.cfgId = cfgId

	s.logic = logic2.GetLogicFactory().Create(logicType)
	s.logic.Init(s)
	s.space = s.logic.CreateSpace()

	s.world.SetContext(s)
	s.world.SetEventListener(s)

	s.startTime = s.timeProvider.NowMs()

	s.entry = csv.Scene.GetEntry(int(cfgId))
	s.width = int(s.entry.Width)
	s.height = int(s.entry.Height)
	s.graph = blockgraph.NewGraph(int(s.entry.Width), int(s.entry.Height))
	s.initAStar()
}

func (s *Scene) initAStar() {
	config := astar.Config{
		GridWidth:  s.width,
		GridHeight: s.height,
	}
	finder, err := astar.New(config)
	if err != nil {
		return
	}
	s.finder = finder
	s.pathContext = newPathContext(s.graph)
}

func (s *Scene) checkBlock(x, y int) bool {
	return s.graph.IsInBlock(x, y)
}

func (s *Scene) Start() {
	s.logic.Start()
}

func (s *Scene) GetNodeService() *service.NodeService {
	return s.ns
}

func (s *Scene) GetSceneId() uint64 {
	return s.sceneId
}

func (s *Scene) GetCfgId() int32 {
	return s.cfgId
}

func (s *Scene) GetWorld() entity.IWorld {
	return s.world
}

func (s *Scene) GetSpace() define.ISpace {
	return s.space
}

func (s *Scene) GetGraph() *blockgraph.Graph {
	return s.graph
}

func (s *Scene) GetEvents() *light.EventCenter {
	return s.events
}

func (s *Scene) GetWorldForFight() common2.IWorld {
	return s.worldForFight
}

func (s *Scene) NowMs() int64 {
	return s.timeProvider.NowMs()
}

func (s *Scene) SetInitData(data define.ISceneInitData) {
	s.initData = data
	s.logic.SetSceneInitData(data)
}

func (s *Scene) GetInitData() define.ISceneInitData {
	return s.initData
}

func (s *Scene) SetFatalErr() {
	s.fatalErr = true
}

func (s *Scene) IsOver() bool {
	return s.fatalErr || s.logic.IsOver()
}

func (s *Scene) Update() {
	s.timeProvider.Update()
	s.logic.Update()
	s.world.Update()
}

func (s *Scene) Destroy() {
	s.destroyed = true
	// 踢出所有玩家(踢掉线)
	s.KickAllPlayers()
	// 删除所有entity
	s.logic.Destroy()

	// 通知scenem
	app.Notify(s.ns, "scenem.remote.freescene", nil, &messages.SMOnSceneEnd{
		SceneId: s.sceneId,
	})
}

func (s *Scene) CheckToken(token int) bool {
	return token == s.token
}

func (s *Scene) CanEnter(token int) bool {
	return token == s.token
}

func (s *Scene) GetPlayer(uid int64) define.IPlayer {
	return s.players[uid]
}

func (s *Scene) GetSceneAge() int64 {
	return s.timeProvider.NowMs() - s.startTime
}

func (s *Scene) GetPlayerNum() int {
	return len(s.players)
}

// OnClientSceneLoadOver 客户端场景载入完毕
func (s *Scene) OnClientSceneLoadOver(uid int64) {
	s.logic.OnClientSceneLoadOver(uid)
}

func (s *Scene) OnPlayerOffline() {
}

func (s *Scene) OnPlayerOnline() {
}

func (s *Scene) tryDestroyPlayerCamera(player define.IPlayer) {
	// 删除camera
	if player.HasCamera() {
		logger := s.ns.GetLogger()

		logger.Infof("destroy camera entity: %v", player.GetCamera())
		s.GetWorld().DestroyEntity(player.GetCamera())
		player.ClearCamera()
	}
}

func (s *Scene) DestroyCamera(uid int64) {
	logger := s.ns.GetLogger()
	logger.Infof("DestroyCamera: %v", uid)

	scene := s
	player := scene.GetPlayer(uid)
	if player == nil {
		return
	}

	s.tryDestroyPlayerCamera(player)
}

func (s *Scene) DestroyPlayerBindEntities(uid int64) {
	logger := s.ns.GetLogger()
	logger.Infof("DestroyPlayerBindEntities: %v", uid)
	scene := s
	player := scene.GetPlayer(uid)
	if player == nil {
		logger.Errorf("DestroyPlayerBindEntities, can not find player %v", uid)
		return
	}

	s.tryDestroyPlayerCamera(player)

	// 删除其他对象
	if player.HasAvatar() {
		scene.GetWorld().DestroyEntity(player.GetAvatar())
		player.ClearAvatar()
	}
}
