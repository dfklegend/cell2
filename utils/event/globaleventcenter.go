package event

import (
	"sync"
	"sync/atomic"

	//"github.com/dfklegend/cell2/utils/common"
	"github.com/dfklegend/cell2/utils/logger"
)

var globalEventCenter = newGlobalEC()

type LocalECList struct {
	localECs sync.Map
	Size     int32
}

func newLocalECList() *LocalECList {
	return &LocalECList{}
}

func (l *LocalECList) Add(ec ILocalEventCenter) {
	l.localECs.Store(ec.GetId(), ec)
	atomic.AddInt32(&l.Size, 1)
}

func (l *LocalECList) Del(ec ILocalEventCenter) {
	l.localECs.Delete(ec.GetId())
	atomic.AddInt32(&l.Size, -1)
}

func (l *LocalECList) Range(f func(k, v interface{}) bool) {
	l.localECs.Range(f)
}

type GlobalEventCenter struct {
	events sync.Map
}

func newGlobalEC() *GlobalEventCenter {
	return &GlobalEventCenter{}
}

func GetGlobalEC() *GlobalEventCenter {
	return globalEventCenter
}

func (g *GlobalEventCenter) getECList(eventName string, createIfMiss bool) *LocalECList {
	v, ok := g.events.Load(eventName)
	var l *LocalECList
	if !ok {
		if !createIfMiss {
			return nil
		}
		l = newLocalECList()
		g.events.Store(eventName, l)
		return l
	}
	return v.(*LocalECList)
}

func (g *GlobalEventCenter) Subscribe(eventName string, child ILocalEventCenter) {
	listeners := g.getECList(eventName, true)
	listeners.Add(child)
}

func (g *GlobalEventCenter) Unsubscribe(eventName string, child ILocalEventCenter) {
	listeners := g.getECList(eventName, false)
	if listeners == nil {
		return
	}
	listeners.Del(child)
}

func (g *GlobalEventCenter) Publish(eventName string, args ...interface{}) {
	listeners := g.getECList(eventName, false)
	if listeners == nil {
		return
	}

	e := &EObj{
		EventName: eventName,
		Args:      args,
	}
	listeners.Range(func(k, v interface{}) bool {
		l := v.(ILocalEventCenter)

		select {
		case l.GetChanEvent() <- e:
			return true
		default:
			// 警告，某个事件队列一直没收取消息
			logger.Log.Warnf("localEventCenter:%v queue full", l.GetId())
			return true
		}
		return true
	})
}

func (g *GlobalEventCenter) DumpInfo() {
	logger.Log.Debugf("GlobalEventCenter info:")
	g.events.Range(func(k, v interface{}) bool {
		en := k.(string)
		el := v.(*LocalECList)
		logger.Log.Debugf("%v size:%v", en, el.Size)
		return true
	})
}
