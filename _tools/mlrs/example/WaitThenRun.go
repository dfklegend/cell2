package example

import (
	"log"

	"mlrs/b3"
	"mlrs/b3/config"
	"mlrs/b3/core"
	"mlrs/b3/factory"
)

func init() {
	factory.RegisterMap().Register("WaitThenRun", &WaitThenRun{})
}

// WaitThenRun WaitThenRun(<waitTime>)
type WaitThenRun struct {
	core.Decorator

	waitTime string
}

func (n *WaitThenRun) Initialize(params *config.BTNodeCfg) {
	n.BaseNode.Initialize(params)
	//TODO 初始化变量
}

func (n *WaitThenRun) OnOpen(tick *core.Tick) {
	log.Println(n.GetName())
}

func (n *WaitThenRun) OnTick(tick *core.Tick) b3.Status {
	return n.GetChild().Execute(tick)
}
