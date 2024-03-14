package service

type CreateFunc func(name string)

// FuncCreator impl IServiceCreator
type FuncCreator struct {
	f CreateFunc
}

func NewFuncCreator(f CreateFunc) *FuncCreator {
	return &FuncCreator{
		f: f,
	}
}

func (c *FuncCreator) Create(name string) {
	c.f(name)
}
