package example

import (
	"log"

	"mlrs/b3"
	"mlrs/b3/config"
	"mlrs/b3/core"
	"mlrs/b3/factory"
)

func init() {
	factory.RegisterMap().Register("GetTime", &GetTime{})
}

// GetTime GetTime
type GetTime struct {
	core.Action
}

func (n *GetTime) Initialize(params *config.BTNodeCfg) {
	n.BaseNode.Initialize(params)
}

func (n *GetTime) OnOpen(tick *core.Tick) {
	log.Println(n.GetName())
}

func (n *GetTime) OnTick(tick *core.Tick) b3.Status {
	return b3.SUCCESS
}
