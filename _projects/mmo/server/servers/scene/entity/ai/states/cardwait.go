package states

import (
	ai2 "mmo/servers/scene/entity/ai"
	fsm2 "mmo/servers/scene/entity/fsm"
)

func init() {
	ai2.GetStateFactory().Register(ai2.CardStateWait, fsm2.NewFuncCreator(func() fsm2.IState {
		return NewCardWait()
	}))
}

type CardWait struct {
	*fsm2.BaseState
}

func NewCardWait() fsm2.IState {
	return &CardWait{
		BaseState: fsm2.NewBaseState(),
	}
}

func (r *CardWait) OnEnter() {
}

func (r *CardWait) OnLeave() {
}

func (r *CardWait) Update() {
}

func (r *CardWait) IsOver() bool {
	return false
}
