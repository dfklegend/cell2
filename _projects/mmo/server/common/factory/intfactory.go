package factory

type IObject interface {
}

type ICreator interface {
	Create(args ...any) IObject
}

type IntFactory struct {
	creators map[int]ICreator
}

func NewIntFactory() *IntFactory {
	return &IntFactory{
		creators: map[int]ICreator{},
	}
}

func (f *IntFactory) Register(t int, creator ICreator) {
	f.creators[t] = creator
}

func (f *IntFactory) RegisterFunc(t int, cb func(args ...any) IObject) {
	f.Register(t, &FuncCreator{
		cb: cb,
	})
}

func (f *IntFactory) Create(t int, args ...any) IObject {
	creator := f.creators[t]
	if creator == nil {
		return nil
	}
	return creator.Create(args...)
}

func (f *IntFactory) Visit(visitor func(t int, creator ICreator)) {
	for k, v := range f.creators {
		visitor(k, v)
	}
}
