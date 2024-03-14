package timer

import (
	"sync"
	"time"

	"github.com/dfklegend/cell2/utils/common"
	l "github.com/dfklegend/cell2/utils/logger"
)

/**
 * timer到期后，被push到channel
 * 便于整合到具体使用线程
 * 此timer性能一般，后续考虑做一个性能更好的
 */

type IdType uint64
type CBFunc func(...interface{})
type QueueChan chan *Obj

var timerMgr = NewTimerMgr()

type ITimerMgr interface {
	After(duration time.Duration, cb CBFunc, args ...interface{}) IdType
	AddTimer(duration time.Duration, cb func(), args ...interface{}) IdType
	Cancel(id IdType)
	GetQueue() QueueChan
}

type Obj struct {
	TimerId  IdType
	Duration time.Duration
	CB       CBFunc
	Args     []interface{}
	Canceled bool
	timer    *time.Timer
}

// NewTimerObj TODO: 可以考虑pool
func NewTimerObj(id IdType,
	duration time.Duration,
	cb CBFunc,
	args []interface{}) *Obj {
	return &Obj{
		TimerId:  id,
		Duration: duration,
		CB:       cb,
		Args:     args,
		Canceled: false,
	}
}

// Mgr 每个环境可以创建自己的timerMgr
type Mgr struct {
	idService *common.SerialIdService64
	queue     QueueChan

	// 用于cancel
	// TimerIdType:*TimerObj
	timers  sync.Map
	running bool
}

func NewTimerMgr() *Mgr {
	return &Mgr{
		idService: common.NewSerialIdService64(),
		queue:     make(QueueChan, 999),
		running:   true,
	}
}

//	用于简单测试
//	用户可以创建自己的
func GetTimerMgr() *Mgr {
	return timerMgr
}

func (m *Mgr) Stop() {
	m.running = false
}

func (m *Mgr) allocId() IdType {
	return IdType(m.idService.AllocId())
}

func (m *Mgr) After(duration time.Duration, cb CBFunc, args ...interface{}) IdType {
	t := NewTimerObj(m.allocId(),
		0, cb, args)

	m.doLater(duration, t)
	m.timers.Store(t.TimerId, t)
	return t.TimerId
}

func (m *Mgr) AddTimer(duration time.Duration, cb CBFunc, args ...interface{}) IdType {
	t := NewTimerObj(m.allocId(),
		duration, cb, args)

	m.doLater(duration, t)
	m.timers.Store(t.TimerId, t)
	return t.TimerId
}

func (m *Mgr) Cancel(id IdType) {
	v, ok := m.timers.Load(id)
	if !ok {
		return
	}
	t := v.(*Obj)
	t.Canceled = true
	if t.timer != nil {
		t.timer.Stop()
	}

	m.timers.Delete(id)
}

func (m *Mgr) doLater(duration time.Duration, t *Obj) {
	t.timer = time.AfterFunc(duration, func() {
		//log.Println("cb")
		if t.Canceled {
			return
		}

		// 停止了
		if !m.running {
			return
		}
		// 底层在新routine执行
		// time.AfterFunc
		m.queue <- t
	})
}

/*
	select {
		case t :<- mgr.GetQueue():
			mgr.Do(t)
	}
*/
// 外部从queue读取后，调用
func (m *Mgr) Do(t *Obj) {
	if t.Canceled {
		return
	}

	//t.CB(t.Args...)
	m.do(t)

	// CB内被cancel掉了
	if t.Canceled {
		return
	}

	if t.Duration > 0 {
		m.doLater(t.Duration, t)
	} else {
		m.timers.Delete(t.TimerId)
	}
}

func (m *Mgr) do(t *Obj) {
	defer func() {
		if err := recover(); err != nil {
			l.E.Errorf("panic in timer.do:%v", err)
			l.E.Errorf(common.GetStackStr())
		}
	}()

	t.CB(t.Args...)
}

func (m *Mgr) GetQueue() QueueChan {
	return m.queue
}
