package example

import (
	"log"
	"reflect"

	"mlrs/b3"
	"mlrs/b3/config"
	"mlrs/b3/core"
	"mlrs/b3/factory"
)

func init() {
	factory.RegisterMap().Register("IsConnect", &IsConnect{})
}

// IsConnect IsConnect
type IsConnect struct {
	core.Condition

	isConnected string
}

func (n *IsConnect) Initialize(params *config.BTNodeCfg) {
	n.BaseNode.Initialize(params)
}

func (n *IsConnect) OnOpen(tick *core.Tick) {
	log.Println(n.GetName())
}

func (n *IsConnect) OnTick(tick *core.Tick) b3.Status {

	login := tick.Blackboard.GetMem("login")
	if reflect.TypeOf(login) != nil {
		return b3.SUCCESS
	}
	return b3.FAILURE
}
