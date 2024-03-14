package example

import (
	"log"
	"os"

	"mlrs/b3"
	"mlrs/b3/config"
	"mlrs/b3/core"
	"mlrs/b3/factory"
)

func init() {
	factory.RegisterMap().Register("Quit", &Quit{})
}

// Quit 退出
type Quit struct {
	core.Action
}

func (n *Quit) Initialize(params *config.BTNodeCfg) {
	n.BaseNode.Initialize(params)
}

func (n *Quit) OnOpen(tick *core.Tick) {
	log.Println(n.GetName())
}

func (n *Quit) OnTick(tick *core.Tick) b3.Status {
	os.Exit(0)
	return b3.SUCCESS
}
