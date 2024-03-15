package attr

type funcWatcher struct {
	cb func(attr IAttr)
}

func NewFuncWatcher(cb func(attr IAttr)) *funcWatcher {
	return &funcWatcher{
		cb: cb,
	}
}

func (f *funcWatcher) OnDirt(attr IAttr) {
	f.cb(attr)
}

type funcWatcher1 func(attr IAttr)

func (f funcWatcher1) OnDirt(attr IAttr) {
	f(attr)
}
