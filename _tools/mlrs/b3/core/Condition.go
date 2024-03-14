package core

import (
	"mlrs/b3"
	"mlrs/b3/config"
)

type ICondition interface {
	IBaseNode
}

type Condition struct {
	BaseNode
	BaseWorker
}

func (node *Condition) Ctor() {
	node.category = b3.CONDITION
}

func (node *Condition) Initialize(params *config.BTNodeCfg) {
	node.BaseNode.Initialize(params)
	//node.BaseNode.IBaseWorker = node
}
