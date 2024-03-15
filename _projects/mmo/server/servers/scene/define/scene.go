package define

import (
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/event/light"

	"mmo/common/entity"
	common2 "mmo/modules/fight/common"
	"mmo/servers/scene/blockgraph"
)

const (
	SceneTypeNormal     = iota + 1 // 永久存在
	SceneTypePlayerInst            // 玩家副本，玩家不存在即删除
)

// IScene
// 场景接口对象
type IScene interface {
	SetInitData(data ISceneInitData)
	GetInitData() ISceneInitData

	IsOver() bool
	PlayerLeave(uid int64)

	GetWorld() entity.IWorld
	GetSpace() ISpace
	GetNodeService() *service.NodeService
	GetEvents() *light.EventCenter
	GetWorldForFight() common2.IWorld

	PushViewMsg(pos Pos, route string, msg any)
	GetPlayer(uid int64) IPlayer

	AddCamera(id entity.EntityID, camera ICamera)
	RemoveCamera(id entity.EntityID)
	OnCameraEnter(camera ICamera)

	DestroyCamera(uid int64)
	DestroyPlayerBindEntities(uid int64)
	// PushViewSnapshot 广播entity快照消息
	PushViewSnapshot(e entity.IEntity)

	// GetSceneAge 副本存活时间
	GetSceneAge() int64
	GetPlayerNum() int
	NowMs() int64

	GetSceneId() uint64
	GetCfgId() int32

	GetEntity(id entity.EntityID) entity.IEntity

	GetGraph() *blockgraph.Graph
	ToGridX(float32) int
	ToGridZ(float32) int

	FindPath(src, tar Pos) *Path
	IsValidPos(tar Pos) bool
	IsInBlock(tar Pos) bool
}

// ISceneInitData
// 有些scene有初始化数据，比如对战阵容之类的
// 不同战斗，在场景启动时，会有一些初始化数据传递过去
type ISceneInitData interface {
}

// ISearcher 搜索器
// 在空间中，搜索指定对象的
type ISearcher interface {
	// Validate
	// 对象是否符合要求
	// 比如 是否是敌人
	Validate(id entity.EntityID, dist float32) bool
	AddCandidate(id entity.EntityID, dist float32)
	// MakeResults 取最终的结果
	MakeResults() []entity.EntityID
}

// ISpace
// 提供空间管理，便于快速搜索
// 比如 分成小区块的对象管理
type ISpace interface {
	AddEntity(id entity.EntityID, pos Pos)
	RemoveEntity(id entity.EntityID)
	// UpdateEntityPos 更新位置
	UpdateEntityPos(id entity.EntityID, pos Pos)

	// SearchCircleTargets 圆形范围寻路
	SearchCircleTargets(pos Pos, radius float32, searcher ISearcher) []entity.EntityID
}

// ICamera
// 扩展思考
// 代表可以看见的对象
// 玩家可以创建多个camera
// 一般来说，玩家会创建一个camera
type ICamera interface {
	// GetCameraId 摄像机id, 主摄像机0
	GetCameraId() int32
	GetFrontId() string
	GetNetId() uint32
}

// IViewStrategy 可见策略的抽象
// 比如: 全场景可见, 视野可见
type IViewStrategy interface {
	PushViewMsg(pos Pos, route string, msg any)
}
