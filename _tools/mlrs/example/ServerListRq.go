package example

import (
	"mlrs/b3"
	"mlrs/b3/config"
	"mlrs/b3/core"
	"mlrs/b3/factory"
)

func init() {
	factory.RegisterMap().Register("ServerListRq", &ServerListRq{})
}

// ServerListRq 服务器列表
type ServerListRq struct {
	core.Action

	url string
}

func (n *ServerListRq) Initialize(params *config.BTNodeCfg) {
	n.BaseNode.Initialize(params)
	//TODO 初始化变量
}

func (n *ServerListRq) OnOpen(tick *core.Tick) {

}

func (n *ServerListRq) OnTick(tick *core.Tick) b3.Status {
	return b3.SUCCESS
}
