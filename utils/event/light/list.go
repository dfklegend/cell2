package light

type ListenerList struct {
	Listeners map[uint64]*EListener
}

func NewListenerList() *ListenerList {
	return &ListenerList{
		Listeners: make(map[uint64]*EListener),
	}
}

func (li *ListenerList) Add(l *EListener) {
	li.Listeners[l.Id] = l
}

func (li *ListenerList) Del(id uint64) {
	delete(li.Listeners, id)
}

func (li *ListenerList) GetSize() int {
	return len(li.Listeners)
}

func (li *ListenerList) FindId(cb CBFunc) uint64 {
	pointer := toPointer(cb)
	for k, v := range li.Listeners {
		if pointer == v.GetCBPointer() {
			return k
		}
	}
	return 0
}

// FindIdWithReceiver
// 对象的Pointer是否会改变？
func (li *ListenerList) FindIdWithReceiver(receiver any, cb CBFunc) uint64 {
	pointer := toPointer(cb)
	for k, v := range li.Listeners {
		if v.Receiver != nil {
			if pointer == v.GetCBPointer() && receiver == v.Receiver {
				return k
			}
		} else {
			if pointer == v.GetCBPointer() {
				return k
			}
		}

	}
	return 0
}
