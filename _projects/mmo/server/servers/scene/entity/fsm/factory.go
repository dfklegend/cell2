package fsm

type StateFactory struct {
	creators map[int]IStateCreator
}

func NewStateFactory() *StateFactory {
	return &StateFactory{
		creators: map[int]IStateCreator{},
	}
}

func (f *StateFactory) Register(state int, creator IStateCreator) {
	f.creators[state] = creator
}

func (f *StateFactory) CreateState(state int) IState {
	creator := f.creators[state]
	if creator == nil {
		return nil
	}
	return creator.Create()
}
