package states

import (
	ai2 "mmo/servers/scene/entity/ai"
	fsm2 "mmo/servers/scene/entity/fsm"
)

func init() {
	ai2.GetStateFactory().Register(ai2.StateInit, fsm2.NewFuncCreator(func() fsm2.IState {
		return NewStateInit()
	}))
}

// 如果非移动
// 每秒有10%概率随机移动

type StateInit struct {
	*fsm2.BaseState
}

func NewStateInit() fsm2.IState {
	return &StateInit{
		BaseState: fsm2.NewBaseState(),
	}
}

func (r *StateInit) OnEnter() {
}

func (r *StateInit) OnLeave() {
}

func (r *StateInit) Update() {
}

func (r *StateInit) IsOver() bool {
	return false
}
