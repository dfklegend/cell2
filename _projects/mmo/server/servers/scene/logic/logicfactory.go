package logic

var (
	logicFactory = newLogicFactory()
)

func GetLogicFactory() ISceneLogicFactory {
	return logicFactory
}

type Factory struct {
	creators map[string]Creator
}

func newLogicFactory() *Factory {
	return &Factory{
		creators: make(map[string]Creator),
	}
}

func (f *Factory) Register(name string, creator Creator) {
	f.creators[name] = creator
}

func (f *Factory) create(name string) ISceneLogic {
	creator := f.creators[name]
	if creator == nil {
		return nil
	}
	return creator()
}

func (f *Factory) Create(name string) ISceneLogic {
	creator := f.creators[name]
	if creator == nil {
		return f.create("empty")
	}
	return creator()
}
