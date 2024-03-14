package example

import (
	"log"

	"mlrs/b3"
	"mlrs/b3/config"
	"mlrs/b3/core"
	"mlrs/b3/factory"
)

func init() {
	factory.RegisterMap().Register("ServerAuth", &ServerAuth{})
}

// ServerAuth 服务器认证
type ServerAuth struct {
	core.Action

	account string

	passwd string
}

func (n *ServerAuth) Initialize(params *config.BTNodeCfg) {
	n.BaseNode.Initialize(params)
	//TODO 初始化变量
}

func (n *ServerAuth) OnOpen(tick *core.Tick) {
	log.Println(n.GetName())
}

func (n *ServerAuth) OnTick(tick *core.Tick) b3.Status {
	return b3.SUCCESS
}
