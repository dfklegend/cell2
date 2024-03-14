package bridge

type FnGet func(args ...any) any

type Center struct {
	fns map[string]FnGet
}

func NewCenter() *Center {
	return &Center{
		fns: map[string]FnGet{},
	}
}

func (c *Center) Register(name string, fn FnGet) {
	c.fns[name] = fn
}

func (c *Center) Get(name string, args ...any) any {
	fn := c.fns[name]
	if fn == nil {
		return nil
	}
	return fn(args...)
}
