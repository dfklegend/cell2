package example

import (
	"log"

	"mlrs/b3"
	"mlrs/b3/config"
	"mlrs/b3/core"
	"mlrs/b3/factory"
)

func init() {
	factory.RegisterMap().Register("RandWaitThenRun", &RandWaitThenRun{})
}

// RandWaitThenRun 随机等待(<min_time> to <max_time>)ms之后执行
type RandWaitThenRun struct {
	core.Decorator

	min_time string

	max_time string
}

func (n *RandWaitThenRun) Initialize(params *config.BTNodeCfg) {
	n.BaseNode.Initialize(params)
	//TODO 初始化变量
}

func (n *RandWaitThenRun) OnOpen(tick *core.Tick) {
	log.Println(n.GetName())
}

func (n *RandWaitThenRun) OnTick(tick *core.Tick) b3.Status {
	return b3.SUCCESS
}
