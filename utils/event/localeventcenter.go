package event

import (
	"sync"

	"github.com/dfklegend/cell2/utils/common"
	"github.com/dfklegend/cell2/utils/logger"
)

var (
	idService    = common.NewSerialIdService64()
	ecIdService  = common.NewSerialIdService64()
	useFakeMutex = false
)

func NewRWMutex() common.IMutex {
	if useFakeMutex {
		return &common.FakeMutex{}
	}
	return &sync.RWMutex{}
}

// --------

type ListenerList struct {
	Global    bool
	RWLock    common.IMutex
	Listeners map[uint64]*EListener
}

func NewListenerList() *ListenerList {
	return &ListenerList{
		RWLock:    NewRWMutex(),
		Listeners: make(map[uint64]*EListener),
	}
}

func (li *ListenerList) Add(l *EListener) {
	li.RWLock.Lock()
	defer li.RWLock.Unlock()
	li.Listeners[l.Id] = l
}

func (li *ListenerList) Del(id uint64) {
	li.RWLock.Lock()
	defer li.RWLock.Unlock()
	delete(li.Listeners, id)
}

func (li *ListenerList) GetSize() int {
	return len(li.Listeners)
}

// --------

func NewEventListener() *EListener {
	return &EListener{
		Id:   idService.AllocId(),
		Args: make([]interface{}, 0),
	}
}

type LocalEventCenter struct {
	id     uint64
	RWLock common.IMutex
	events map[string]*ListenerList
	// 接收global
	chanEvent ChanEvent
	// 是否都使用chan
	// 包括本地事件
	localUseChan bool
	running      bool
}

func NewLocalEventCenter(useChan bool) *LocalEventCenter {
	return &LocalEventCenter{
		id:           ecIdService.AllocId(),
		events:       make(map[string]*ListenerList),
		chanEvent:    make(ChanEvent, 999),
		localUseChan: useChan,
		RWLock:       NewRWMutex(),
		running:      true,
	}
}

func (le *LocalEventCenter) Clear() {
	// 需要取消所有全局的事件注册
	le.running = false
	le.RWLock.Lock()
	defer le.RWLock.Unlock()

	for en, el := range le.events {
		if !el.Global {
			continue
		}

		GetGlobalEC().Unsubscribe(en, le)
	}

	le.events = make(map[string]*ListenerList)
}

func (le *LocalEventCenter) SetLocalUseChan(v bool) {
	le.localUseChan = v
}

func (le *LocalEventCenter) GetId() uint64 {
	return le.id
}

func (le *LocalEventCenter) getListener(eventName string, createIfMiss bool) *ListenerList {
	le.RWLock.RLock()
	v, ok := le.events[eventName]
	le.RWLock.RUnlock()
	if !ok {
		if !createIfMiss {
			return nil
		}
		v = NewListenerList()

		le.RWLock.Lock()
		if le.running {
			le.events[eventName] = v
		}
		le.RWLock.Unlock()
	}
	return v
}

func (le *LocalEventCenter) Subscribe(eventName string, cb CBFunc, args ...interface{}) uint64 {
	if !le.running {
		return 0
	}
	list := le.getListener(eventName, true)
	l := NewEventListener()
	l.Args = args
	l.CB = cb
	list.Add(l)
	return l.Id
}

// GSubscribe 先注册自己
// 向Global注册
func (le *LocalEventCenter) GSubscribe(eventName string, cb CBFunc, args ...interface{}) uint64 {
	if !le.running {
		return 0
	}
	list := le.getListener(eventName, true)
	id := le.Subscribe(eventName, cb, args...)
	if !list.Global {
		// 向global注册
		GetGlobalEC().Subscribe(eventName, le)
		list.Global = true
	}
	return id
}

func (le *LocalEventCenter) Unsubscribe(eventName string, id uint64) {
	list := le.getListener(eventName, false)
	if list == nil {
		return
	}
	list.Del(id)

	// 发现自己空了，global节点取消
	if list.Global && list.GetSize() == 0 {
		GetGlobalEC().Unsubscribe(eventName, le)
		list.Global = false
	}
}

func (le *LocalEventCenter) GUnsubscribe(eventName string, id uint64) {
	le.Unsubscribe(eventName, id)
}

func (le *LocalEventCenter) GetChanEvent() ChanEvent {
	return le.chanEvent
}

// DoEvent call by chan receiver
func (le *LocalEventCenter) DoEvent(e *EObj) {
	// 派发
	le.dispatch(e.EventName, e.Args...)
}

func (le *LocalEventCenter) Publish(eventName string, args ...interface{}) {
	if le.localUseChan {
		e := &EObj{
			EventName: eventName,
			Args:      args,
		}
		le.chanEvent <- e
	} else {
		le.dispatch(eventName, args...)
	}
}

func (le *LocalEventCenter) dispatch(eventName string, args ...interface{}) {
	ll := le.getListener(eventName, false)
	if ll == nil {
		return
	}
	listeners := ll.Listeners

	ll.RWLock.RLock()
	for _, v := range listeners {
		finalArgs := append(v.Args, args...)
		v.CB(finalArgs...)
	}
	ll.RWLock.RUnlock()
}

func (le *LocalEventCenter) DumpInfo(eventName string) {
	ll := le.getListener(eventName, false)
	if ll == nil {
		return
	}

	ll.RWLock.RLock()
	logger.Log.Debugf("%v size:%v", eventName, ll.GetSize())
	ll.RWLock.RUnlock()
}
