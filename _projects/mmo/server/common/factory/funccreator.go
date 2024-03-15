package factory

type FuncCreator struct {
	cb func(args ...any) IObject
}

func (c *FuncCreator) Create(args ...any) IObject {
	return c.cb(args...)
}
