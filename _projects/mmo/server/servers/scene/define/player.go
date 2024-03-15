package define

import (
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/event/light"

	"mmo/common/define"
	"mmo/common/entity"
	mymsg "mmo/messages"
	"mmo/modules/systems"
)

// IPlayer 玩家对象
type IPlayer interface {
	GetNodeService() *service.NodeService

	GetId() int64
	GetFrontId() string
	GetNetId() uint32
	GetLogicId() string
	SetOnline(bool)
	IsOnline() bool

	//SetEntityId(id entity.EntityID)
	//GetEntityId() entity.EntityID

	SetAvatar(id entity.EntityID)
	GetAvatar() entity.EntityID
	GetAvatarEntity() entity.IEntity
	HasAvatar() bool
	ClearAvatar()

	SetCamera(id entity.EntityID)
	HasCamera() bool
	GetCamera() entity.EntityID
	ClearCamera()

	ChangeState(state ScenePlayerState)
	GetState() ScenePlayerState
	IsState(state ScenePlayerState) bool

	EnterWorld(newPlayer bool, info *mymsg.PlayerInfo, switchLine bool)
	MakeSaveInfo() *mymsg.PlayerInfo

	IsDirt() bool
	SetDirt()
	ClearDirt()

	UpdateScene(scene IScene, cfgId int32, sceneId uint64)
	CanSwitchLine() bool
	ExitPointReqChangeScene(tarCfgId int, pos define.Pos)
	PushMsg(route string, msg interface{})
	Kick()

	GetSystem(name string) systems.ISystem
	GetEvents() *light.EventCenter
}
