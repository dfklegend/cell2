package core

import (
	"mlrs/b3"
	. "mlrs/b3/config"
)

type IAction interface {
	IBaseNode
}

type Action struct {
	BaseNode
	BaseWorker
}

func (node *Action) Ctor() {
	node.category = b3.ACTION
}
func (node *Action) Initialize(params *BTNodeCfg) {

	//node.id = b3.CreateUUID()
	node.BaseNode.Initialize(params)
	//node.BaseNode.IBaseWorker = node
	node.parameters = make(map[string]interface{})
	node.properties = make(map[string]interface{})
}
