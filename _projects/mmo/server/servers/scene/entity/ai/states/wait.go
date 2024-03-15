package states

import (
	ai2 "mmo/servers/scene/entity/ai"
	fsm2 "mmo/servers/scene/entity/fsm"
)

func init() {
	ai2.GetStateFactory().Register(ai2.StateWait, fsm2.NewFuncCreator(func() fsm2.IState {
		return NewWait()
	}))
}

type Wait struct {
	*fsm2.BaseState
}

func NewWait() fsm2.IState {
	return &Wait{
		BaseState: fsm2.NewBaseState(),
	}
}

func (r *Wait) OnEnter() {
}

func (r *Wait) OnLeave() {
}

func (r *Wait) Update() {
}

func (r *Wait) IsOver() bool {
	return false
}
