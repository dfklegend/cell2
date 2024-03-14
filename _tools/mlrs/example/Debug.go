package example

import (
	"log"

	"mlrs/b3"
	"mlrs/b3/config"
	"mlrs/b3/core"
	"mlrs/b3/factory"
)

func init() {
	factory.RegisterMap().Register("Debug", &Debug{})
}

// Debug Debug.log(<msg>)
type Debug struct {
	core.Action

	msg string
}

func (n *Debug) Initialize(params *config.BTNodeCfg) {
	n.BaseNode.Initialize(params)
	n.msg = params.GetPropertyAsString("msg")
}

func (n *Debug) OnOpen(tick *core.Tick) {
	log.Println(n.msg)
}

func (n *Debug) OnTick(tick *core.Tick) b3.Status {
	return b3.SUCCESS
}
