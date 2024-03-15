package fsm

type BaseCtrl struct {
	ICtrl
	factory   IFactory
	stateType int
	state     IState

	ctx    IContext
	cbCtrl CtrlFunc
}

func NewBaseCtrl() *BaseCtrl {
	return &BaseCtrl{}
}

func (c *BaseCtrl) Init(factory IFactory, f CtrlFunc, ctx IContext) {
	c.factory = factory
	c.cbCtrl = f
	c.ctx = ctx
}

func (c *BaseCtrl) GetStateType() int {
	return c.stateType
}

func (c *BaseCtrl) GetState() IState {
	return c.state
}

func (c *BaseCtrl) GetContext() IContext {
	return c.ctx
}

func (c *BaseCtrl) Update() {
	if c.state != nil && !c.state.IsOver() {
		c.state.Update()
	}

	if c.cbCtrl != nil {
		c.cbCtrl(c)
	}
}

func (c *BaseCtrl) ChangeState(state int) {
	newState := c.factory.CreateState(state)
	if newState == nil {
		return
	}
	if c.state != nil {
		c.state.OnLeave()
	}

	c.stateType = state
	c.state = newState
	newState.Init(c)
	newState.OnEnter()
}
