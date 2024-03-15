package components

import (
	"github.com/dfklegend/cell2/utils/common"

	define2 "mmo/servers/scene/define"
	"mmo/servers/scene/entity/define"
	"mmo/servers/scene/space/utils"
)

// ExitComponent
// 场景出口
// 切换到另外的场景
// 定期搜索范围内玩家角色，传送到目标位置
type ExitComponent struct {
	*BaseSceneComponent

	tarCfgId int32
	tarPos   define2.Pos

	radius    float32
	nextCheck int64
	tran      *Transform
}

func NewExitComponent(tarCfgId int32, tarPos define2.Pos, radius float32) *ExitComponent {
	return &ExitComponent{
		BaseSceneComponent: NewBaseSceneComponent(),
		radius:             radius,
		tarCfgId:           tarCfgId,
		tarPos:             tarPos,
	}
}

func (t *ExitComponent) OnPrepare() {
	t.BaseSceneComponent.OnPrepare()
}

func (t *ExitComponent) OnStart() {
	t.tran = t.GetOwner().GetComponent(define.Transform).(*Transform)
	t.nextCheck = common.NowMs()
}

func (t *ExitComponent) Update() {
	t.tryTransportPlayers()
}

func (t *ExitComponent) LateUpdate() {
}

func (t *ExitComponent) OnDestroy() {
}

func (t *ExitComponent) tryTransportPlayers() {
	now := common.NowMs()
	if now < t.nextCheck {
		return
	}
	t.nextCheck = now + 1000
	t.transportPlayers()
}

func (t *ExitComponent) transportPlayers() {
	tars := utils.FindPlayers(t.GetScene().GetSpace(), t.GetOwner(), t.radius)
	if len(tars) == 0 {
		return
	}

	logger := t.GetScene().GetNodeService().GetLogger()
	// 每一个玩家对象
	// 要求他发起切换场景请求
	scene := t.GetScene()
	for _, v := range tars {
		e := scene.GetEntity(v)
		if e == nil {
			// log error
			return
		}
		p := e.GetComponent(define.Player).(*PlayerComponent)
		player := p.GetPlayer()
		// 发送切换场景请求
		//player.GetAvatar()
		logger.Infof("transport player: %v", player.GetId())
		player.ExitPointReqChangeScene(int(t.tarCfgId), define2.Pos{})
	}
}
