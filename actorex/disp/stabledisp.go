package disp

// 创建一个固定携程调度，而不是动态创建
type stableDisp struct {
	chanTask ChanTask
}

func NewStableDisp() *stableDisp {
	return &stableDisp{
		chanTask: make(ChanTask, 999),
	}
}

func (d *stableDisp) Schedule(fn func()) {
	d.chanTask <- fn
}

func (d *stableDisp) Throughput() int {
	return 999
}

func (d *stableDisp) Start() {
	go d.run()
}

func (d *stableDisp) run() {
	for {
		select {
		case fn := <-d.chanTask:
			fn()
		}
	}
}
