package example

import (
	"log"

	"mlrs/b3"
	"mlrs/b3/config"
	"mlrs/b3/core"
	"mlrs/b3/factory"
)

func init() {
	factory.RegisterMap().Register("LoginRq", &LoginRq{})
}

// LoginRq LoginRq
type LoginRq struct {
	core.Action
}

func (n *LoginRq) Initialize(params *config.BTNodeCfg) {
	n.BaseNode.Initialize(params)
}

func (n *LoginRq) OnOpen(tick *core.Tick) {
	log.Println(n.GetName())
	tick.Blackboard.SetMem("login", true)
}

func (n *LoginRq) OnTick(tick *core.Tick) b3.Status {
	return b3.SUCCESS
}
