package example

import (
	"log"

	"mlrs/b3"
	"mlrs/b3/config"
	"mlrs/b3/core"
	"mlrs/b3/factory"
)

func init() {
	factory.RegisterMap().Register("SelectServer", &SelectServer{})
}

// SelectServer 选择服务器
type SelectServer struct {
	core.Action

	server string
}

func (n *SelectServer) Initialize(params *config.BTNodeCfg) {
	n.BaseNode.Initialize(params)
	//TODO 初始化变量
}

func (n *SelectServer) OnOpen(tick *core.Tick) {
	log.Println(n.GetName())
}

func (n *SelectServer) OnTick(tick *core.Tick) b3.Status {
	return b3.SUCCESS
}
