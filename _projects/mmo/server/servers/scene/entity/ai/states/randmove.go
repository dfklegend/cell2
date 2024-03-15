package states

import (
	"math/rand"

	"github.com/dfklegend/cell2/utils/common"

	define2 "mmo/servers/scene/define"
	ai2 "mmo/servers/scene/entity/ai"
	"mmo/servers/scene/entity/ai/ctrl"
	fsm2 "mmo/servers/scene/entity/fsm"
)

func init() {
	ai2.GetStateFactory().Register(ai2.StateRandMove, fsm2.NewFuncCreator(func() fsm2.IState {
		return NewRandMove()
	}))
}

// 如果非移动
// 每秒有10%概率随机移动

type RandMove struct {
	*fsm2.BaseState

	nextCheck int64
}

func NewRandMove() fsm2.IState {
	return &RandMove{
		BaseState: fsm2.NewBaseState(),
	}
}

func (r *RandMove) OnEnter() {
}

func (r *RandMove) OnLeave() {
}

func (r *RandMove) Update() {
	if common.NowMs() < r.nextCheck {
		return
	}
	r.nextCheck = common.NowMs() + 1000

	ai := r.GetCtrl().GetContext().(*ctrl.AICtrl)
	if ai.GetMove().IsMoving() {
		return
	}
	if rand.Float32() > 0.1 {
		return
	}

	tran := ai.GetTran()
	tar := tran.GetPos()

	tar.X += (rand.Float32() - 0.5) * 10
	tar.Z += (rand.Float32() - 0.5) * 10

	tar.X = clampPos(tar.X)
	tar.Z = clampPos(tar.Z)

	ai.GetMove().MoveTo(tar)
}

func (r *RandMove) IsOver() bool {
	return false
}

func clampPos(one float32) float32 {
	size := define2.MaxWidth
	if one < -size {
		return -size
	}
	if one > size {
		return size
	}
	return one
}
