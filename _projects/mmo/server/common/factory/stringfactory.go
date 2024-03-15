package factory

type StringFactory struct {
	creators map[string]ICreator
}

func NewStringFactory() *StringFactory {
	return &StringFactory{
		creators: map[string]ICreator{},
	}
}

func (f *StringFactory) Register(t string, creator ICreator) {
	f.creators[t] = creator
}

func (f *StringFactory) RegisterFunc(t string, cb func(args ...any) IObject) {
	f.Register(t, &FuncCreator{
		cb: cb,
	})
}

func (f *StringFactory) Create(t string, args ...any) IObject {
	creator := f.creators[t]
	if creator == nil {
		return nil
	}
	return creator.Create(args...)
}

func (f *StringFactory) Visit(visitor func(name string, creator ICreator)) {
	for k, v := range f.creators {
		visitor(k, v)
	}
}
