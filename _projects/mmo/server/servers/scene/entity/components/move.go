package components

import (
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/common"

	"mmo/messages/cproto"
	define3 "mmo/servers/scene/define"
	"mmo/servers/scene/entity/define"
)

// MoveComponent 负责移动
type MoveComponent struct {
	*BaseSceneComponent
	tran *Transform

	moving        bool
	tar           define3.Pos
	speed         float32
	nextChangeTar int64

	lastUpdate int64
}

func NewMoveComponent() *MoveComponent {
	return &MoveComponent{
		BaseSceneComponent: NewBaseSceneComponent(),
		moving:             false,
		speed:              1,
	}
}

func (m *MoveComponent) OnPrepare() {
	m.BaseSceneComponent.OnPrepare()
}

func (m *MoveComponent) OnStart() {
	m.tran = m.GetOwner().GetComponent(define.Transform).(*Transform)
	m.lastUpdate = common.NowMs()
}

func (m *MoveComponent) Update() {
	m.updateMoving()
}

func (m *MoveComponent) OnDestroy() {
}

func (m *MoveComponent) MoveTo(tar define3.Pos) {
	m.tar = tar

	if !m.moving {
		m.StartMove()
	}

	m.scene.PushViewMsg(define3.Pos{}, "moveto", &cproto.MoveTo{
		Id: int32(m.GetOwner().GetId()),
		Tar: &cproto.Vector3{
			X: tar.X,
			Z: tar.Z,
		},
	})
}

func (m *MoveComponent) StopMove() {
	m.EndMove()

	pos := m.tran.GetPos()
	m.scene.PushViewMsg(define3.Pos{}, "stopmove", &cproto.MoveTo{
		Id: int32(m.GetOwner().GetId()),
		Tar: &cproto.Vector3{
			X: pos.X,
			Z: pos.Z,
		},
	})
}

func (m *MoveComponent) updateMoving() {
	if !m.moving {
		return
	}
	now := common.NowMs()
	time := float32((now - m.lastUpdate) / 1000)
	if time == 0 {
		return
	}
	m.lastUpdate = now

	dir := m.tar.Sub(m.tran.pos)
	dist := dir.Magnitude()
	dir = dir.Normalized()

	step := m.speed * time
	if step >= dist {
		m.setPos(m.tar)
		m.EndMove()
		return
	}

	pos := m.tran.pos.Add(dir.Mul(step))
	m.setPos(pos)
}

func (m *MoveComponent) setPos(pos define3.Pos) {
	m.tran.SetPos(pos)
}

func (m *MoveComponent) GetTar() define3.Pos {
	return m.tar
}

func (m *MoveComponent) IsMoving() bool {
	return m.moving
}

func (m *MoveComponent) PushMoveTo(ns *service.NodeService, camera define3.ICamera) {
	tar := m.tar
	app.PushMessageById(ns, camera.GetFrontId(), camera.GetNetId(), "moveto", &cproto.MoveTo{
		Id: int32(m.GetOwner().GetId()),
		Tar: &cproto.Vector3{
			X: tar.X,
			Z: tar.Z,
		}})
}

func (m *MoveComponent) StartMove() {
	m.moving = true
	m.lastUpdate = common.NowMs()
}

func (m *MoveComponent) EndMove() {
	m.moving = false
}

func clampPos(one float32) float32 {
	size := define3.MaxWidth
	if one < -size {
		return -size
	}
	if one > size {
		return size
	}
	return one
}
