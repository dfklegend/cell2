package scenem

import (
	"time"

	"github.com/dfklegend/cell2/utils/timer"

	"mmo/common/config"
	"mmo/common/define"
)

type PublicScene struct {
	CfgId   int32
	ReqNum  int32
	Spawned int32
}

// PublicScenes 负责维护公共场景
// 时机合适后，刷场景
// 发现场景少了，会补
type PublicScenes struct {
	mgr     *SceneServiceMgr
	spawned int

	scenes map[int32]*PublicScene

	timerUpdate timer.IdType
}

func NewPublicScenes() *PublicScenes {
	return &PublicScenes{
		scenes: make(map[int32]*PublicScene),
	}
}

func (p *PublicScenes) Init(mgr *SceneServiceMgr) {
	p.mgr = mgr

	if config.PerfTest {
		p.addPublicScene(100, 1)
		p.addPublicScene(101, 50)
		p.addPublicScene(102, 50)
	}

	if config.EnablePublicScene {
		p.addPublicScene(100, 1)
		p.addPublicScene(101, 5)
		p.addPublicScene(102, 5)
	}
}

func (p *PublicScenes) addPublicScene(cfgId int32, reqNum int32) {
	scene := p.scenes[cfgId]
	if scene != nil {
		return
	}
	scene = &PublicScene{
		CfgId:  cfgId,
		ReqNum: reqNum,
	}

	p.scenes[cfgId] = scene
}

func (p *PublicScenes) getScene(cfgId int32) *PublicScene {
	return p.scenes[cfgId]
}

func (p *PublicScenes) Start() {
	if p.timerUpdate != 0 {
		return
	}
	p.timerUpdate = p.mgr.GetNodeService().GetRunService().GetTimerMgr().AddTimer(
		time.Second, func(args ...any) {
			p.Update()
		})
}

func (p *PublicScenes) Update() {
	for _, v := range p.scenes {
		p.trySpawnScene(v)
	}
}

func (p *PublicScenes) SpawnOne() bool {
	return p.mgr.SpawnScene(int32(define.SceneNormal))
}

// 保证刷新固定场景
func (p *PublicScenes) trySpawnScene(scene *PublicScene) {
	scene.Spawned = int32(p.mgr.world.GetLineNum(scene.CfgId))
	if scene.Spawned >= scene.ReqNum {
		return
	}

	if p.mgr.SpawnScene(scene.CfgId) {
		scene.Spawned++
	}
}
