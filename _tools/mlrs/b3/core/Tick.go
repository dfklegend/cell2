package core

import (
	_ "fmt"
)

type Tick struct {
	tree *BehaviorTree

	debug interface{}

	target interface{}

	Blackboard *Blackboard

	_openNodes []IBaseNode

	_openSubtreeNodes []*SubTree

	_nodeCount int
}

func NewTick() *Tick {
	tick := &Tick{}
	tick.Initialize()
	return tick
}

func (tick *Tick) Initialize() {
	// set by BehaviorTree
	tick.tree = nil
	tick.debug = nil
	tick.target = nil
	tick.Blackboard = nil

	// updated during the tick signal
	tick._openNodes = nil
	tick._openSubtreeNodes = nil
	tick._nodeCount = 0
}

func (tick *Tick) GetTree() *BehaviorTree {
	return tick.tree
}

func (tick *Tick) _enterNode(node IBaseNode) {
	tick._nodeCount++
	tick._openNodes = append(tick._openNodes, node)

	// TODO: call debug here
}

func (tick *Tick) _openNode(node *BaseNode) {
	// TODO: call debug here
}

func (tick *Tick) _tickNode(node *BaseNode) {
	// TODO: call debug here
	//fmt.Println("Tick _tickNode :", this.debug, " id:", node.GetID(), node.GetTitle())
}

func (tick *Tick) _closeNode(node *BaseNode) {
	// TODO: call debug here

	ulen := len(tick._openNodes)
	if ulen > 0 {
		tick._openNodes = tick._openNodes[:ulen-1]
	}

}

func (tick *Tick) pushSubtreeNode(node *SubTree) {
	tick._openSubtreeNodes = append(tick._openSubtreeNodes, node)
}
func (tick *Tick) popSubtreeNode() {
	len := len(tick._openSubtreeNodes)
	if len > 0 {
		tick._openSubtreeNodes = tick._openSubtreeNodes[:len-1]
	}
}

func (tick *Tick) GetLastSubTree() *SubTree {
	ulen := len(tick._openSubtreeNodes)
	if ulen > 0 {
		return tick._openSubtreeNodes[ulen-1]
	}
	return nil
}

func (tick *Tick) _exitNode(node *BaseNode) {
	// TODO: call debug here
}

func (tick *Tick) GetTarget() interface{} {
	return tick.target
}
