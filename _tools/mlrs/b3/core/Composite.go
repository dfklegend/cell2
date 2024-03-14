package core

import (
	"mlrs/b3"
	"mlrs/b3/config"
)

type IComposite interface {
	IBaseNode
	GetChildCount() int
	GetChild(index int) IBaseNode
	AddChild(child IBaseNode)
}

type Composite struct {
	BaseNode
	BaseWorker

	children []IBaseNode
}

func (node *Composite) Ctor() {
	node.category = b3.COMPOSITE
}

func (node *Composite) Initialize(params *config.BTNodeCfg) {
	node.BaseNode.Initialize(params)
	//node.BaseNode.IBaseWorker = node
	node.children = make([]IBaseNode, 0)
	//fmt.Println("Composite Initialize")
}

func (node *Composite) GetChildCount() int {
	return len(node.children)
}

func (node *Composite) GetChild(index int) IBaseNode {
	return node.children[index]
}

func (node *Composite) AddChild(child IBaseNode) {
	node.children = append(node.children, child)
}

func (node *Composite) tick(tick *Tick) b3.Status {
	return b3.ERROR
}
