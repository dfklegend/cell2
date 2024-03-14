package core

import (
	"mlrs/b3"
	"mlrs/b3/config"
)

type IDecorator interface {
	IBaseNode
	SetChild(child IBaseNode)
	GetChild() IBaseNode
}

type Decorator struct {
	BaseNode
	BaseWorker
	child IBaseNode
}

func (node *Decorator) Ctor() {
	node.category = b3.DECORATOR
}

func (node *Decorator) Initialize(params *config.BTNodeCfg) {
	node.BaseNode.Initialize(params)
	//node.BaseNode.IBaseWorker = node
}

func (node *Decorator) GetChild() IBaseNode {
	return node.child
}

func (node *Decorator) SetChild(child IBaseNode) {
	node.child = child
}
