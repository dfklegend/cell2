package fsm

// IState state可以决定自己是否结束，但是不要决定要切到那个状态
// 可以增加state的重用性(AI)
// 在CtrlFunc中来决定切换状态
type IState interface {
	Init(ICtrl)
	GetCtrl() ICtrl

	OnEnter()
	OnLeave()

	Update()
	IsOver() bool
}

type IStateCreator interface {
	Create() IState
}

type IFactory interface {
	Register(stateType int, creator IStateCreator)
	CreateState(stateType int) IState
}

// CtrlFunc 可以进行状态主动切换
type CtrlFunc func(ctrl ICtrl)

type IContext interface {
}

type ICtrl interface {
	Init(factory IFactory, f CtrlFunc, ctx IContext)

	GetStateType() int
	GetState() IState
	GetContext() IContext
	ChangeState(state int)

	Update()
}
