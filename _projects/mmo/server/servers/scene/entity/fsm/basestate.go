package fsm

type BaseState struct {
	ctrl ICtrl
}

func NewBaseState() *BaseState {
	return &BaseState{}
}

func (s *BaseState) Init(ctrl ICtrl) {
	s.ctrl = ctrl
}

func (s *BaseState) GetCtrl() ICtrl {
	return s.ctrl
}
