package sche

import (
	"sync/atomic"

	"github.com/dfklegend/cell2/utils/common"
	l "github.com/dfklegend/cell2/utils/logger"
)

/*
	sche包
	提供调度器，调度器用来在runservice中提供统一调度
	MultiSelector提供了方便的可扩展的select
*/

//
var (
	TheScheMgr     *Mgr = NewScheMgr()
	DefaultScheMgr *Mgr = TheScheMgr
	// 此机制虽然可以避免一种死锁，但是造成客户端消息顺序不一致，先关闭
	selfBlockDefend = false
)

// 	最简单的匿名函数
type CBFunc func()

// --------------------------

type RunTask struct {
	id uint32
	cb CBFunc
}

// 	--------------------------
type RunTaskIdService struct {
	nextId uint32
}

func (s *RunTaskIdService) AllocId() uint32 {
	return atomic.AddUint32(&s.nextId, 1)
}

var (
	runtaskidservice RunTaskIdService = RunTaskIdService{
		nextId: 1,
	}
)

// --------------------------

type ISche interface {
	AddTask(cb CBFunc) *RunTask
}

type Sche struct {
	chanTask  chan *RunTask
	chanClose chan int
}

const (
	QueueSize = 999
)

func NewSche() *Sche {
	return &Sche{
		chanTask:  make(chan *RunTask, QueueSize),
		chanClose: make(chan int, 1),
	}
}

func (s *Sche) Stop() {
	close(s.chanTask)
	close(s.chanClose)
}

// Post
// 如果利用waterfall.sche来执行
// 本地队列满了，同时又Post到本Sche，应该会死锁
func (s *Sche) Post(cb CBFunc) *RunTask {
	defer func() {
		if err := recover(); err != nil {
			l.E.Errorf("panic in sche.Post:%v", err)
			l.E.Errorf(common.GetStackStr())
		}
	}()

	t := &RunTask{
		id: runtaskidservice.AllocId(),
		cb: cb,
	}

	if selfBlockDefend && len(s.chanTask) >= QueueSize-10 {
		// 快满了，防止自身routine post 的死锁
		go func() {
			defer func() {
				if err := recover(); err != nil {
					l.E.Errorf("panic in selfBlockDefend sche.Post:%v", err)
					l.E.Errorf(common.GetStackStr())
				}
			}()

			s.chanTask <- t
		}()
		return t
	}

	s.chanTask <- t
	return t
}

// GetChanTask 集成在其他流程中
/*
	select {
	case t := <- sche.GetChanTask():
		sche.DoTask(t)
	}
*/
func (s *Sche) GetChanTask() <-chan *RunTask {
	return s.chanTask
}

func (s *Sche) doTask(t *RunTask) {
	defer func() {
		if err := recover(); err != nil {
			l.E.Errorf("panic in sche.doTask:%v", err)
			l.E.Errorf(common.GetStackStr())
		}
	}()

	t.cb()
}

func (s *Sche) DoTask(t *RunTask) {
	s.doTask(t)
}

// Handler 便于测试，自己go sche.Handler
func (s *Sche) Handler() {
	for true {
		select {
		case t := <-s.chanTask:
			if t != nil {
				s.DoTask(t)
			}
		case <-s.chanClose:
			return
		}
	}
}
