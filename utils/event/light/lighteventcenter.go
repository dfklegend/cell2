package light

import (
	"log"

	"github.com/dfklegend/cell2/utils/common"
	"github.com/dfklegend/cell2/utils/logger"
)

/*
 * 非线程安全
 * Sub/Unsubscribe, 如果是对象，请传入对象实例
 * SubscribeWithReceiver
 * (TODO: 如果反射可以分析出func是对象的method，并获取receiver，可以简化接口)
 */

type CBFunc func(args ...interface{})

type EventCenter struct {
	idService *common.SerialIdService64
	events    map[string]*ListenerList

	running bool
}

func NewEventCenter() *EventCenter {
	return &EventCenter{
		idService: common.NewSerialIdService64(),
		events:    make(map[string]*ListenerList),
		running:   true,
	}
}

func (c *EventCenter) allocId() uint64 {
	return c.idService.AllocId()
}

func (c *EventCenter) Clear() {
	c.running = false
	c.events = make(map[string]*ListenerList)
}

func (c *EventCenter) getListener(eventName string, createIfMiss bool) *ListenerList {
	v, ok := c.events[eventName]
	if !ok {
		if !createIfMiss {
			return nil
		}
		v = NewListenerList()

		if c.running {
			c.events[eventName] = v
		}
	}
	return v
}

func (c *EventCenter) subscribeCheck(checkCB bool, eventName string, cb CBFunc, args ...interface{}) uint64 {
	if !c.running {
		return 0
	}
	list := c.getListener(eventName, true)

	if checkCB {
		id := list.FindId(cb)
		if id != 0 {
			// already subscribe
			logger.L.Warnf("already register %v with this CB", eventName)
			return 0
		}
	}

	l := NewEventListener(c.allocId())
	l.Args = args
	l.CB = cb
	list.Add(l)

	return l.Id
}

func (c *EventCenter) SubscribeNoCheck(eventName string, cb CBFunc, args ...interface{}) uint64 {
	return c.subscribeCheck(false, eventName, cb, args...)
}

func (c *EventCenter) Subscribe(eventName string, cb CBFunc, args ...interface{}) uint64 {
	return c.subscribeCheck(true, eventName, cb, args...)
}

// SubscribeWithReceiver receiver必须是指针对象
func (c *EventCenter) SubscribeWithReceiver(eventName string, receiver any, cb CBFunc, args ...interface{}) uint64 {
	if !c.running {
		return 0
	}
	list := c.getListener(eventName, true)

	id := list.FindIdWithReceiver(receiver, cb)
	if id != 0 {
		// already subscribe
		log.Printf("already register %v with this receiver and CB", eventName)
		return 0
	}

	l := NewEventListener(c.allocId())
	l.Args = args
	l.CB = cb

	l.Receiver = receiver
	list.Add(l)

	return l.Id
}

func (c *EventCenter) UnsubscribeById(eventName string, id uint64) {
	list := c.getListener(eventName, false)
	if list == nil {
		return
	}
	list.Del(id)
}

func (c *EventCenter) Unsubscribe(eventName string, cb CBFunc) {
	list := c.getListener(eventName, false)
	if list == nil {
		return
	}
	id := list.FindId(cb)
	if id == 0 {
		return
	}
	list.Del(id)
}

func (c *EventCenter) UnsubscribeWithReceiver(eventName string, receiver any, cb CBFunc) {
	list := c.getListener(eventName, false)
	if list == nil {
		return
	}
	id := list.FindIdWithReceiver(receiver, cb)
	if id == 0 {
		return
	}
	list.Del(id)
}

func (c *EventCenter) Publish(eventName string, args ...interface{}) {
	c.dispatch(eventName, args...)
}

func (c *EventCenter) dispatch(eventName string, args ...interface{}) {
	ll := c.getListener(eventName, false)
	if ll == nil {
		return
	}
	listeners := ll.Listeners

	for _, v := range listeners {
		if len(v.Args) > 0 {
			finalArgs := append(v.Args, args...)
			v.CB(finalArgs...)
		} else {
			v.CB(args...)
		}

	}
}

// GetSubscribeNum 订阅事件的个数
func (c *EventCenter) GetSubscribeNum(eventName string) int {
	ll := c.getListener(eventName, false)
	if ll == nil {
		return 0
	}
	return ll.GetSize()
}

func (c *EventCenter) HasSubscribers(eventName string) bool {
	return c.GetSubscribeNum(eventName) > 0
}

func (c *EventCenter) DumpInfo(eventName string) {
	ll := c.getListener(eventName, false)
	if ll == nil {
		return
	}

	logger.L.Infof("%v size:%v", eventName, ll.GetSize())
}
