package actions

import (
	"fmt"

	"mlrs/b3"
	. "mlrs/b3/config"
	. "mlrs/b3/core"
)

type Log struct {
	Action
	info string
}

func (node *Log) Initialize(setting *BTNodeCfg) {
	node.Action.Initialize(setting)
	node.info = setting.GetPropertyAsString("info")
}

func (node *Log) OnTick(tick *Tick) b3.Status {
	fmt.Println("log:", node.info)
	return b3.SUCCESS
}
