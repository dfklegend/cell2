package mailbox

import (
	"runtime"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/dfklegend/cell2/actorex/queue/mpsc"
	"github.com/dfklegend/cell2/utils/common"
	l "github.com/dfklegend/cell2/utils/logger"
)

/*
尽量平滑的处理消息
一次执行超过了maxProcessCost，将会跳出循环，并停止调度1ms
(时间可以被交给同service的其他处理者，避免太多阻塞)
*/
type queue interface {
	Push(interface{})
	Pop() interface{}
}

const (
	idle int32 = iota
	running

	MaxMsgNumToSmooth = 100000
)

type SmoothFrameMailbox struct {
	userMailbox     queue
	systemMailbox   *mpsc.Queue
	schedulerStatus int32
	userMessages    int32
	sysMessages     int32
	suspended       int32

	smoothPaused int32

	invoker     actor.MessageInvoker
	dispatcher  actor.Dispatcher
	middlewares []actor.MailboxMiddleware

	// smoothpause
	// ns
	maxProcessCost int64
}

func (m *SmoothFrameMailbox) PostUserMessage(message interface{}) {
	// is it a raw batch message?
	if batch, ok := message.(actor.MessageBatch); ok {
		messages := batch.GetMessages()

		for _, msg := range messages {
			m.PostUserMessage(msg)
		}
	}

	// is it an envelope batch message?
	if env, ok := message.(actor.MessageEnvelope); ok {
		if batch, ok := env.Message.(actor.MessageBatch); ok {
			messages := batch.GetMessages()

			for _, msg := range messages {
				m.PostUserMessage(msg)
			}
		}
	}

	// normal messages
	for _, ms := range m.middlewares {
		ms.MessagePosted(message)
	}
	m.userMailbox.Push(message)
	atomic.AddInt32(&m.userMessages, 1)
	m.schedule()
}

func (m *SmoothFrameMailbox) PostSystemMessage(message interface{}) {
	for _, ms := range m.middlewares {
		ms.MessagePosted(message)
	}
	m.systemMailbox.Push(message)
	atomic.AddInt32(&m.sysMessages, 1)
	m.schedule()
}

func (m *SmoothFrameMailbox) RegisterHandlers(invoker actor.MessageInvoker, dispatcher actor.Dispatcher) {
	m.invoker = invoker
	m.dispatcher = dispatcher
}

func (m *SmoothFrameMailbox) schedule() {
	if atomic.LoadInt32(&m.smoothPaused) == 1 {
		return
	}
	if atomic.CompareAndSwapInt32(&m.schedulerStatus, idle, running) {
		m.dispatcher.Schedule(m.processMessages)
	}
}

func (m *SmoothFrameMailbox) processMessages() {

	m.run()

	// set mailbox to idle
	atomic.StoreInt32(&m.schedulerStatus, idle)
	sys := atomic.LoadInt32(&m.sysMessages)
	user := atomic.LoadInt32(&m.userMessages)
	smoothPaused := atomic.LoadInt32(&m.smoothPaused)

	// check if there are still messages to process (sent after the message loop ended)
	if sys > 0 || (atomic.LoadInt32(&m.suspended) == 0 && user > 0 && smoothPaused == 0) {
		// 又来了新消息, 重新调度
		m.schedule()
		return
	}

	//
	if sys == 0 && user == 0 {
		for _, ms := range m.middlewares {
			ms.MailboxEmpty()
		}
	}
}

func (m *SmoothFrameMailbox) beginSmoothPause() {
	// 暂停一下
	if !atomic.CompareAndSwapInt32(&m.smoothPaused, 0, 1) {
		return
	}
	go func() {
		// 释放一点时间
		time.Sleep(time.Millisecond)
		atomic.CompareAndSwapInt32(&m.smoothPaused, 1, 0)

		//l.Log.Infof("[ACTOR] %v %v end sleep", unsafe.Pointer(m), NowNano())
		m.schedule()
	}()
}

func (m *SmoothFrameMailbox) run() {
	var msg interface{}

	defer func() {
		if r := recover(); r != nil {
			//l.Log.Infof("[ACTOR] Recovering", log.Object("actor", m.invoker), log.Object("reason", r), log.Stack())
			//l.E.Errorf("[ACTOR] Recovering", log.Object("actor", m.invoker), log.Object("reason", r), log.Stack())
			l.Log.Infof("[ACTOR] Recovering", m.invoker, r)
			l.E.Errorf("[ACTOR] Recovering", m.invoker, r)

			l.E.Errorf("error: %v", r)
			stack := common.GetStackStr()
			l.E.Errorf(stack)

			m.invoker.EscalateFailure(r, msg)
		}
	}()

	i, t := 0, m.dispatcher.Throughput()

	beginTime := NowNano()
	processed := 0
	for {
		if i > t {
			i = 0
			//runtime.Gosched()
		}

		// 避免消息太多
		cost := NowNano() - beginTime
		msgNum := m.userMessages
		if msgNum < MaxMsgNumToSmooth {
			// 超过20ms
			if cost > m.maxProcessCost {
				//l.Log.Infof("[ACTOR] %v %v beginSmoothPause cost: %v msgs: %v processed: %v",
				//	unsafe.Pointer(m), NowNano(), cost, msgNum, processed)
				m.beginSmoothPause()
				return
			}
		} else {
			// to many msgs just sched it
			if cost > m.maxProcessCost {
				l.Log.Infof("[ACTOR] %v %v Gosched cost: %v msgs: %v processed: %v",
					unsafe.Pointer(m), NowNano(), cost, msgNum, processed)
				processed = 0
				beginTime = NowNano()
				runtime.Gosched()
			}
		}

		i++

		// keep processing system messages until queue is empty
		if msg = m.systemMailbox.Pop(); msg != nil {
			atomic.AddInt32(&m.sysMessages, -1)
			switch msg.(type) {
			case *actor.SuspendMailbox:
				atomic.StoreInt32(&m.suspended, 1)
			case *actor.ResumeMailbox:
				atomic.StoreInt32(&m.suspended, 0)
			default:
				m.invoker.InvokeSystemMessage(msg)
			}
			for _, ms := range m.middlewares {
				ms.MessageReceived(msg)
			}
			continue
		}

		// didn't process a system message, so break until we are resumed
		if atomic.LoadInt32(&m.suspended) == 1 {
			return
		}

		if msg = m.userMailbox.Pop(); msg != nil {
			atomic.AddInt32(&m.userMessages, -1)
			m.invoker.InvokeUserMessage(msg)
			processed++
			for _, ms := range m.middlewares {
				ms.MessageReceived(msg)
			}
		} else {
			//if cost > 5000000 {
			//	l.Log.Infof("[ACTOR] %v %v process over cost: %v msgs: %v processed: %v",
			//		unsafe.Pointer(m), NowNano(), cost, msgNum, processed)
			//}
			return
		}
	}
}

func (m *SmoothFrameMailbox) Start() {
	for _, ms := range m.middlewares {
		ms.MailboxStarted()
	}
}

func (m *SmoothFrameMailbox) UserMessageCount() int {
	return int(atomic.LoadInt32(&m.userMessages))
}
