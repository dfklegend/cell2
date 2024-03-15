package scenelogics

import (
	"fmt"
	"math/rand"

	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/service"
	l "github.com/dfklegend/cell2/utils/logger"

	"mmo/common/config"
	"mmo/common/entity"
	"mmo/messages/cproto"
	"mmo/modules/csv"
	common2 "mmo/modules/fight/common"
	"mmo/modules/scenecfg"
	define3 "mmo/servers/scene/define"
	builder2 "mmo/servers/scene/entity/builder"
	components2 "mmo/servers/scene/entity/components"
	define2 "mmo/servers/scene/entity/define"
	logic2 "mmo/servers/scene/logic"
	"mmo/servers/scene/space/factory"
)

const (
	IdleDestroyTime = 5 * 60 * 1000
)

func init() {
	logic2.GetLogicFactory().Register("standard", func() logic2.ISceneLogic {
		return newStandard()
	})
}

/*
	Standard
	普通场景
		允许玩家随机加入
		加入后，会自动创建对应的unitAvatar
		根据场景配置创建unitMonster
*/
type Standard struct {
	scene define3.IScene

	//sceneImpl *sceneservice.Scene
	ns *service.NodeService

	maxMonster int32
	monsterNum int32

	timeNoPlayer int64
	permanent    bool

	cfg *scenecfg.SceneCfg
}

func newStandard() *Standard {
	return &Standard{
		maxMonster: 0,
	}
}

func (s *Standard) Init(scene define3.IScene) {
	s.scene = scene
	s.ns = s.scene.GetNodeService()
	s.initFromCfg()
}

func (s *Standard) initFromCfg() {
	cfgId := s.scene.GetCfgId()
	s.cfg = scenecfg.GetMgr().GetItem(int(cfgId))
	entry := csv.Scene.GetEntry(int(cfgId))
	if entry == nil {
		return
	}
	s.maxMonster = entry.MaxMonsterNum
	if config.PerfTest {
		s.maxMonster = 100
	}
	//s.maxMonster = 0
	if entry.Type == define3.SceneTypeNormal {
		s.permanent = true
	}

}

func (s *Standard) CreateSpace() define3.ISpace {
	return factory.CreateZoneSpace()
}

func (s *Standard) PlayerEnter(player define3.IPlayer) {
	scene := s.scene

	if player.HasAvatar() {
		s.ns.GetLogger().Warnf("PlayerEnter already has avatar: %v", player.GetId())
		return
	}

	world := scene.GetWorld()
	eAvatar := builder2.CreateAvatarEntity(world, &builder2.PlayerInfo{
		Player: player,
	})
	player.SetAvatar(eAvatar.GetId())

	s.timeNoPlayer = 0
}

func (s *Standard) PlayerLeave(uid int64) {
	s.scene.DestroyPlayerBindEntities(uid)
}

func (s *Standard) OnClientSceneLoadOver(uid int64) {
	scene := s.scene
	player := scene.GetPlayer(uid)
	if player == nil {
		return
	}

	if player.HasCamera() {
		return
	}

	world := scene.GetWorld()

	// create camera
	e := builder2.CreateCameraEntity(world, player.GetAvatar(), player.GetFrontId(), player.GetNetId())
	player.SetCamera(e.GetId())

	l.L.Infof("OnClientSceneLoadOver: %v, create camera: %v", uid, e.GetId())
	// 下发avatarEnter
	player.PushMsg("avatarenter", &cproto.AvatarEnter{
		Id: int32(player.GetAvatar()),
	})

	player.PushMsg("battlelog", &cproto.BattleLog{
		Log: fmt.Sprintf("进入场景: %v, %v", scene.GetNodeService().Name, scene.GetSceneId()),
	})
}

func (s *Standard) Start() {
	s.createEntities()
}

func (s *Standard) IsOver() bool {
	if s.permanent {
		return false
	}
	if s.scene.GetSceneAge() < 60*1000 {
		return false
	}
	if s.scene.GetPlayerNum() == 0 && s.timeNoPlayer == 0 {
		s.timeNoPlayer = s.scene.NowMs()
	}
	if s.scene.GetPlayerNum() == 0 && s.scene.NowMs() >= s.timeNoPlayer+IdleDestroyTime {
		return true
	}
	return false
}

func (s *Standard) Update() {
	s.trySpawnMonsters()
}

func (s *Standard) Destroy() {
}

func (s *Standard) createEntities() {
	for i := int32(0); i < s.maxMonster; i++ {
		s.createMonster()
	}
	s.createExits()
}

func (s *Standard) createExits() {
	cfg := s.cfg
	if cfg == nil {
		return
	}
	if cfg.Exits == nil {
		return
	}
	world := s.scene.GetWorld()
	for _, v := range cfg.Exits {
		builder2.CreateExitEntity(world, &builder2.ExitInfo{
			Pos: define3.Pos{
				X: v.X,
				Z: v.Z,
			},
			TarCfgId: v.To,
			Radius:   v.Radius,
		})
	}
}

func (s *Standard) trySpawnMonsters() {
	for i := s.monsterNum; i < s.maxMonster; i++ {
		s.createMonster()
	}
}

func (s *Standard) randSelectMonster() string {
	cfg := s.cfg
	if cfg == nil || cfg.Monsters == nil || len(cfg.Monsters) == 0 {
		return ""
	}

	index := rand.Intn(len(cfg.Monsters))
	return cfg.Monsters[index].Id
}

func (s *Standard) createMonster() {
	s.createMonsterById(s.randSelectMonster())
}

func (s *Standard) createMonsterById(cfgId string) {
	md := csv.Monster.GetEntry(cfgId)
	if md == nil {
		return
	}

	side := 2
	if config.MonsterRandSide {
		side = rand.Intn(100)
	}

	info := &builder2.MonsterInfo{
		Name:  md.Name,
		CfgId: cfgId,
		Pos: define3.Pos{
			X: (rand.Float32() - 0.5) * 2 * define3.MaxWidth,
			Z: (rand.Float32() - 0.5) * 2 * define3.MaxWidth,
		},
		Side:  side,
		Level: md.Level,
	}
	builder2.CreateMonsterEntity(s.scene.GetWorld(), info)
}

func (s *Standard) OnPreCameraEnter(camera define3.ICamera) {
}

func (s *Standard) OnPostCameraEnter(camera define3.ICamera) {
}

func (s *Standard) tryAddToSpace(e entity.IEntity) {
	c := e.GetComponent(define2.Transform)
	if c == nil {
		return
	}
	tran := c.(*components2.Transform)
	if tran == nil {
		return
	}
	s.scene.GetSpace().AddEntity(e.GetId(), tran.GetPos())
}

func (s *Standard) OnAddEntity(e entity.IEntity) {
	scene := s.scene

	s.tryAddToSpace(e)
	bu := e.GetComponent(define2.BaseUnit).(*components2.BaseUnit)
	if bu == nil {
		return
	}

	scene.PushViewSnapshot(e)

	switch bu.GetUnitType() {
	case define3.UnitMonster:
		s.monsterNum++
		return
	case define3.UnitCamera:
		s.AddCamera(e)
		s.onCameraEntityEnter(e)
	}
}

func (s *Standard) PushSnapshot(camera define3.ICamera, e entity.IEntity) {
	bu := e.GetComponent(define2.BaseUnit).(*components2.BaseUnit)
	if bu == nil {
		return
	}
	switch bu.GetUnitType() {
	case define3.UnitTest, define3.UnitAvatar:
		s.pushTestSnapshot(camera, e)
	case define3.UnitMonster:
		s.pushMonsterSnapshot(camera, e)
	case define3.UnitExit:
		s.pushExitSnapshot(camera, e)
	}
}

func (s *Standard) pushTestSnapshot(camera define3.ICamera, e entity.IEntity) {
	tran := e.GetComponent(define2.Transform).(*components2.Transform)
	unit := e.GetComponent(define2.BaseUnit).(*components2.BaseUnit)

	pos := &cproto.Vector3{}
	pos.X = tran.GetPos().X
	pos.Z = tran.GetPos().Z

	app.PushMessageById(s.ns, camera.GetFrontId(), camera.GetNetId(), "testsnapshot", &cproto.TestSnapshot{
		Id:    int32(e.GetId()),
		Name:  unit.Name,
		Pos:   pos,
		HP:    int32(unit.GetChar().GetHP()),
		HPMax: int32(unit.GetChar().GetIntValue(common2.HPMax)),
		Side:  int32(unit.GetChar().GetSide()),
	})
}

func (s *Standard) pushMonsterSnapshot(camera define3.ICamera, e entity.IEntity) {
	tran := e.GetComponent(define2.Transform).(*components2.Transform)
	unit := e.GetComponent(define2.BaseUnit).(*components2.BaseUnit)
	monster := e.GetComponent(define2.MonsterComponent).(*components2.MonsterComponent)

	pos := &cproto.Vector3{}
	pos.X = tran.GetPos().X
	pos.Z = tran.GetPos().Z

	app.PushMessageById(s.ns, camera.GetFrontId(), camera.GetNetId(), "monstersnapshot", &cproto.MonsterSnapshot{
		Id:    int32(e.GetId()),
		Name:  unit.Name,
		Pos:   pos,
		HP:    int32(unit.GetChar().GetHP()),
		HPMax: int32(unit.GetChar().GetIntValue(common2.HPMax)),
		Side:  int32(unit.GetChar().GetSide()),
		CfgId: monster.GetCfgId(),
	})
}

func (s *Standard) pushExitSnapshot(camera define3.ICamera, e entity.IEntity) {
	tran := e.GetComponent(define2.Transform).(*components2.Transform)

	pos := &cproto.Vector3{}
	pos.X = tran.GetPos().X
	pos.Z = tran.GetPos().Z

	app.PushMessageById(s.ns, camera.GetFrontId(), camera.GetNetId(), "exitsnapshot", &cproto.ExitSnapshot{
		Id:  int32(e.GetId()),
		Pos: pos,
	})
}

func (s *Standard) AddCamera(e entity.IEntity) {
	camera := e.GetComponent(define2.Camera).(*components2.Camera)
	if camera == nil {
		return
	}
	s.scene.AddCamera(e.GetId(), camera)
}

func (s *Standard) RemoveCamera(e entity.IEntity) {
	camera := e.GetComponent(define2.Camera).(*components2.Camera)
	if camera == nil {
		return
	}
	s.scene.RemoveCamera(e.GetId())
}

func (s *Standard) onCameraEntityEnter(e entity.IEntity) {
	camera := e.GetComponent(define2.Camera).(*components2.Camera)
	if camera == nil {
		return
	}

	s.scene.OnCameraEnter(camera)
}

func (s *Standard) OnDestroyEntity(e entity.IEntity) {
	s.tryRemoveFromSpace(e)
	bu := e.GetComponent(define2.BaseUnit).(*components2.BaseUnit)
	if bu == nil {
		return
	}

	switch bu.GetUnitType() {
	case define3.UnitMonster:
		s.monsterNum--
	case define3.UnitCamera:
		s.RemoveCamera(e)
	}

	s.scene.PushViewMsg(define3.Pos{}, "unitleave", &cproto.UnitLeave{
		Id: int32(e.GetId()),
	})
}

func (s *Standard) tryRemoveFromSpace(e entity.IEntity) {
	c := e.GetComponent(define2.Transform)
	if c == nil {
		return
	}
	tran := c.(*components2.Transform)
	if tran == nil {
		return
	}
	s.scene.GetSpace().RemoveEntity(e.GetId())
}

func (s *Standard) SetSceneInitData(d any) {
}
