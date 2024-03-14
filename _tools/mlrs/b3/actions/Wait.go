package actions

import (
	"time"

	"mlrs/b3"
	. "mlrs/b3/config"
	. "mlrs/b3/core"
)

type Wait struct {
	Action
	endTime int64
}

func (w *Wait) Initialize(setting *BTNodeCfg) {
	w.Action.Initialize(setting)
	w.endTime = setting.GetPropertyAsInt64("milliseconds")
}

func (w *Wait) OnOpen(tick *Tick) {
	var startTime int64 = time.Now().UnixNano() / 1000000
	tick.Blackboard.Set("startTime", startTime, tick.GetTree().GetID(), w.GetID())
}

func (w *Wait) OnTick(tick *Tick) b3.Status {
	var currTime int64 = time.Now().UnixNano() / 1000000
	var startTime = tick.Blackboard.GetInt64("startTime", tick.GetTree().GetID(), w.GetID())
	//fmt.Println("wait:",w.GetTitle(),tick.GetLastSubTree(),"=>", currTime-startTime)
	if currTime-startTime > w.endTime {
		return b3.SUCCESS
	}

	return b3.RUNNING
}
