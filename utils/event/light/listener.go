package light

type EListener struct {
	Id uint64
	// 自带参数
	Args []interface{}
	CB   CBFunc

	Receiver any
}

func NewEventListener(id uint64) *EListener {
	return &EListener{
		Id:   id,
		Args: make([]interface{}, 0),
	}
}

func (l *EListener) GetCBPointer() uintptr {
	return toPointer(l.CB)
}
