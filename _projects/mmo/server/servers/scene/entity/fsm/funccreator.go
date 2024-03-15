package fsm

type FuncCreator struct {
	cb func() IState
}

func NewFuncCreator(cb func() IState) IStateCreator {
	return &FuncCreator{
		cb: cb,
	}
}

func (c *FuncCreator) Create() IState {
	return c.cb()
}
