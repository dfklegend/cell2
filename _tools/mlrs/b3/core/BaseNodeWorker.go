package core

import (
	"fmt"

	"mlrs/b3"
)

type IBaseWorker interface {
	OnEnter(tick *Tick)
	OnOpen(tick *Tick)
	OnTick(tick *Tick) b3.Status
	OnClose(tick *Tick)
	OnExit(tick *Tick)
}
type BaseWorker struct {
}

func (bw *BaseWorker) OnEnter(tick *Tick) {

}

func (bw *BaseWorker) OnOpen(tick *Tick) {

}

func (bw *BaseWorker) OnTick(tick *Tick) b3.Status {
	fmt.Println("tick BaseWorker")
	return b3.ERROR
}

func (bw *BaseWorker) OnClose(tick *Tick) {

}

func (bw *BaseWorker) OnExit(tick *Tick) {

}
