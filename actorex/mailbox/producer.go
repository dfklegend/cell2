package mailbox

import (
	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/actorex/queue/goring"
	"github.com/dfklegend/cell2/actorex/queue/mpsc"
)

type unboundedMailboxQueue struct {
	userMailbox *goring.Queue
}

func (q *unboundedMailboxQueue) Push(m interface{}) {
	q.userMailbox.Push(m)
}

func (q *unboundedMailboxQueue) Pop() interface{} {
	m, o := q.userMailbox.Pop()
	if o {
		return m
	}

	return nil
}

//	Producer
// MaxProcessCost: 最大帧处理时间(ms)
func Producer(maxProcessCostMs int64, mailboxStats ...actor.MailboxMiddleware) actor.MailboxProducer {
	return func() actor.Mailbox {
		q := &unboundedMailboxQueue{
			userMailbox: goring.New(10),
		}
		if maxProcessCostMs == 0 {
			// default
			maxProcessCostMs = 10
		}
		return &SmoothFrameMailbox{
			systemMailbox:  mpsc.New(),
			userMailbox:    q,
			middlewares:    mailboxStats,
			maxProcessCost: maxProcessCostMs * 1000000,
		}
	}
}
