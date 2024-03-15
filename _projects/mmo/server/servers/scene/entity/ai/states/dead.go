package states

import (
	ai2 "mmo/servers/scene/entity/ai"
	fsm2 "mmo/servers/scene/entity/fsm"
)

func init() {
	ai2.GetStateFactory().Register(ai2.StateDead, fsm2.NewFuncCreator(func() fsm2.IState {
		return NewDead()
	}))
}

type Dead struct {
	*fsm2.BaseState
}

func NewDead() fsm2.IState {
	return &Dead{
		BaseState: fsm2.NewBaseState(),
	}
}

func (r *Dead) OnEnter() {
}

func (r *Dead) OnLeave() {
}

func (r *Dead) Update() {
}

func (r *Dead) IsOver() bool {
	return false
}
