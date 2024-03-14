package example

import (
	"log"
	"math/rand"
	"time"

	"mlrs/b3"
	"mlrs/b3/config"
	"mlrs/b3/core"
	"mlrs/b3/factory"
)

func init() {
	factory.RegisterMap().Register("RandWait", &RandWait{})
}

// RandWait RandWait <timemini> to <timemax>
type RandWait struct {
	core.Action

	timemini int64

	timemax int64

	wait int64
}

func (n *RandWait) Initialize(params *config.BTNodeCfg) {
	n.BaseNode.Initialize(params)
	//TODO 初始化变量
	n.timemini = params.GetPropertyAsInt64("timemini")
	n.timemax = params.GetPropertyAsInt64("timemax")

}

func (n *RandWait) OnOpen(tick *core.Tick) {
	tick.Blackboard.Set("time", time.Now().UnixMilli(), n.GetID(), n.GetID())
	n.wait = rand.Int63n(n.timemax-n.timemini) + n.timemini
	log.Println(n.GetName())
}

func (n *RandWait) OnTick(tick *core.Tick) b3.Status {
	now := time.Now().UnixMilli()

	diff := now - tick.Blackboard.GetInt64("time", n.GetID(), n.GetID())

	if diff > n.wait {
		return b3.SUCCESS
	}

	return b3.RUNNING
}
