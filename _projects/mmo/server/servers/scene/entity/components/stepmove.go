package components

import (
	"math"

	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/common"

	"mmo/messages/cproto"
	"mmo/modules/utils"
	define3 "mmo/servers/scene/define"
	"mmo/servers/scene/entity/define"
)

const (
	Idle     = iota
	Walking  // 移动中
	WalkWait // 移动等待
)

const (
	WalkCost float32 = 0.5
	WaitCost float32 = 0.25
)

const (
	moveDetailLog = false
)

// StepMoveComponent 负责移动
// 一格一格移动
type StepMoveComponent struct {
	*BaseSceneComponent
	tran *Transform

	//moving        bool
	tar      define3.Pos
	speed    float32
	passTime float32

	curPosIndex int
	path        *define3.Path

	lastUpdate int64
	state      int

	// 玩家操作，最少走一步
	leastOneStepMode bool
	stopNextStep     bool
}

func NewStepMoveComponent(leastOneStepMode bool) *StepMoveComponent {
	return &StepMoveComponent{
		BaseSceneComponent: NewBaseSceneComponent(),
		state:              Idle,
		speed:              1,
		leastOneStepMode:   leastOneStepMode,
	}
}

func (m *StepMoveComponent) OnPrepare() {
	m.BaseSceneComponent.OnPrepare()
}

func (m *StepMoveComponent) OnStart() {
	m.tran = m.GetOwner().GetComponent(define.Transform).(*Transform)
	m.lastUpdate = common.NowMs()
}

func (m *StepMoveComponent) Update() {

	if !m.IsMoving() {
		return
	}

	now := common.NowMs()
	delta := float32(now-m.lastUpdate) / 1000.0
	m.lastUpdate = now

	switch m.state {
	case Walking:
		m.updateWalking(delta)
	case WalkWait:
		m.updateWait(delta)
	}
}

func (m *StepMoveComponent) OnDestroy() {
}

func (m *StepMoveComponent) setState(state int) {
	m.state = state
}

func (m *StepMoveComponent) isState(state int) bool {
	return state == m.state
}

func (m *StepMoveComponent) MoveTo(tar define3.Pos) {
	tar = utils.MakeGridPos(tar)
	if !m.scene.IsValidPos(tar) {
		return
	}

	logger := m.scene.GetNodeService().GetLogger()
	if moveDetailLog {
		pos := m.tran.GetPos()
		logger.Infof("MoveTo: (%v, %v) -> (%v, %v)", pos.X, pos.Z, tar.X, tar.Z)
	}

	m.tran.ClearBlockPos()
	path := m.scene.FindPath(m.tran.GetPos(), tar)
	m.tran.ResetBlockPos()
	if path == nil {
		return
	}

	m.tar = tar
	m.curPosIndex = 0
	m.path = path

	if moveDetailLog {
		logger.Infof("getPath: %v, %v", len(path.Points), path)
	}

	if !m.IsMoving() {
		m.StartMove()
	} else {
		m.ContinueMove()
	}
}

func (m *StepMoveComponent) StopMove() {
	if !m.IsMoving() {
		return
	}

	if m.leastOneStepMode {
		if m.isState(Walking) {
			m.stopNextStep = true
		} else {
			m.EndMove()
		}
	} else {
		m.EndMove()
	}
}

func (m *StepMoveComponent) updateWalking(delta float32) {
	if m.isReachTar() {
		m.EndMove()
		return
	}

	// 判断如果目标被阻挡了，停下来
	tar := m.getNextStepTar()
	if m.scene.IsInBlock(tar) {
		if moveDetailLog {
			m.scene.GetNodeService().GetLogger().Infof("move blocked")
		}
		m.EndMove()
		return
	}

	m.passTime += delta
	// 每秒走一格
	if m.passTime < WalkCost {
		return
	}
	m.walkOneStep()
	m.passTime -= WalkCost

	// end it
	if m.leastOneStepMode && m.stopNextStep {
		m.EndMove()
		m.stopNextStep = false
		return
	}

	if m.isReachTar() {
		m.EndMove()
		return
	}

	m.setState(WalkWait)
}

func (m *StepMoveComponent) updateWait(delta float32) {
	m.passTime += delta
	// 每秒走一格
	if m.passTime < WaitCost {
		return
	}
	m.passTime -= WaitCost
	m.setState(Walking)
}

func (m *StepMoveComponent) isReachTar() bool {
	if m.curPosIndex >= len(m.path.Points) {
		return true
	}
	off := m.tar.Sub(m.tran.GetPos())
	if math.Abs(float64(off.X))+math.Abs(float64(off.Z)) < 0.1 {
		return true
	}
	return false
}

func (m *StepMoveComponent) hasNextStep() bool {
	return m.curPosIndex <= len(m.path.Points)
}

func (m *StepMoveComponent) getNextStepTar() define3.Pos {
	return m.path.Points[m.curPosIndex]
}

func (m *StepMoveComponent) walkOneStep() {
	if m.curPosIndex >= len(m.path.Points) {
		return
	}

	newPos := m.path.Points[m.curPosIndex]
	m.curPosIndex++
	m.setPos(newPos)

	if moveDetailLog {
		logger := m.scene.GetNodeService().GetLogger()
		logger.Infof("oneStep: %v, %v", newPos.X, newPos.Z)
	}

	m.scene.PushViewMsg(define3.Pos{}, "moveto", &cproto.MoveTo{
		Id: int32(m.GetOwner().GetId()),
		Tar: &cproto.Vector3{
			X: newPos.X,
			Z: newPos.Z,
		},
	})
}

func (m *StepMoveComponent) setPos(pos define3.Pos) {
	m.tran.SetPos(pos)
}

func (m *StepMoveComponent) GetTar() define3.Pos {
	return m.tar
}

func (m *StepMoveComponent) IsMoving() bool {
	return m.state == Walking || m.state == WalkWait
}

func (m *StepMoveComponent) PushMoveTo(ns *service.NodeService, camera define3.ICamera) {
	tar := m.tar
	app.PushMessageById(ns, camera.GetFrontId(), camera.GetNetId(), "moveto", &cproto.MoveTo{
		Id: int32(m.GetOwner().GetId()),
		Tar: &cproto.Vector3{
			X: tar.X,
			Z: tar.Z,
		}})
}

func (m *StepMoveComponent) StartMove() {
	m.setState(Walking)
	m.lastUpdate = common.NowMs()
	m.passTime = 0
}

func (m *StepMoveComponent) ContinueMove() {
	m.passTime = 0
	if m.leastOneStepMode && m.stopNextStep {
		m.stopNextStep = false
	}
}

func (m *StepMoveComponent) EndMove() {
	if moveDetailLog {
		logger := m.scene.GetNodeService().GetLogger()
		pos := m.tran.GetPos()
		logger.Infof("endMove: %v, %v", pos.X, pos.Z)
	}
	m.setState(Idle)
}
