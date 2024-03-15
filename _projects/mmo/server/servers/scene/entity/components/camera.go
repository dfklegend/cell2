package components

import (
	"mmo/common/entity"
	define2 "mmo/servers/scene/entity/define"
)

// Camera 可见性
// 等扩展了可见性，可以增加Transform控件
type Camera struct {
	*BaseSceneComponent

	OwnerId   entity.EntityID
	ownerTran *Transform

	tran    *Transform
	FrontId string
	NetId   uint32
}

func NewCamera(ownerId entity.EntityID, frontId string, netId uint32) *Camera {
	return &Camera{
		BaseSceneComponent: NewBaseSceneComponent(),
		OwnerId:            ownerId,
		FrontId:            frontId,
		NetId:              netId,
	}
}

func (t *Camera) OnPrepare() {
	t.BaseSceneComponent.OnPrepare()

	//t.tran = t.GetOwner().GetComponent(define2.Transform).(*Transform)
}

func (t *Camera) initOwner() {
	if t.OwnerId == 0 {
		return
	}
	owner := t.GetOwner().GetWorld().GetEntity(t.OwnerId)
	if owner == nil {
		return
	}
	t.ownerTran = owner.GetComponent(define2.Transform).(*Transform)
}

func (t *Camera) OnStart() {
}

func (t *Camera) Update() {
}

func (t *Camera) LateUpdate() {
	if t.ownerTran == nil {
		return
	}
	//t.tran.SetPos(t.ownerTran.GetPos())
}

func (t *Camera) OnDestroy() {
}

func (t *Camera) GetCameraId() int32 {
	return 0
}

func (t *Camera) GetFrontId() string {
	return t.FrontId
}

func (t *Camera) GetNetId() uint32 {
	return t.NetId
}
