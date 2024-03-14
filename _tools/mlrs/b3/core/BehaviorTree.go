package core

import (
	"fmt"

	"mlrs/b3"
	"mlrs/b3/config"
)

type BehaviorTree struct {
	id string

	title string

	description string

	properties map[string]interface{}

	root IBaseNode

	debug interface{}

	dumpInfo *config.BTTreeCfg
}

func NewBeTree() *BehaviorTree {
	tree := &BehaviorTree{}
	tree.Initialize()
	return tree
}

func (bt *BehaviorTree) Initialize() {
	bt.id = b3.CreateUUID()
	bt.title = "The behavior tree"
	bt.description = "Default description"
	bt.properties = make(map[string]interface{})
	bt.root = nil
	bt.debug = nil
}

func (bt *BehaviorTree) GetID() string {
	return bt.id
}

func (bt *BehaviorTree) GetTitile() string {
	return bt.title
}

func (bt *BehaviorTree) SetDebug(debug interface{}) {
	bt.debug = debug
}

func (bt *BehaviorTree) GetRoot() IBaseNode {
	return bt.root
}

func (bt *BehaviorTree) Load(data *config.BTTreeCfg, maps *b3.RegisterStructMaps, extMaps *b3.RegisterStructMaps) {
	bt.title = data.Title             //|| bt.title;
	bt.description = data.Description // || bt.description;
	bt.properties = data.Properties   // || bt.properties;
	bt.dumpInfo = data
	nodes := make(map[string]IBaseNode)

	// Create the node list (without connection between them)

	for id, s := range data.Nodes {
		spec := &s
		var node IBaseNode

		if spec.Category == "tree" {
			node = new(SubTree)
		} else {
			if extMaps != nil && extMaps.CheckElem(spec.Name) {
				// Look for the name in custom nodes
				if tnode, err := extMaps.New(spec.Name); err == nil {
					node = tnode.(IBaseNode)
				}
			} else {
				if tnode, err2 := maps.New(spec.Name); err2 == nil {
					node = tnode.(IBaseNode)
				} else {
					//fmt.Println("new ", spec.Name, " err:", err2)
				}
			}
		}

		if node == nil {
			// Invalid node name
			panic("BehaviorTree.load: Invalid node name:" + spec.Name + ",title:" + spec.Title)

		}

		node.Ctor()
		node.Initialize(spec)
		node.SetBaseNodeWorker(node.(IBaseWorker))
		nodes[id] = node
	}

	// Connect the nodes
	for id, spec := range data.Nodes {
		node := nodes[id]

		if node.GetCategory() == b3.COMPOSITE && spec.Children != nil {
			for i := 0; i < len(spec.Children); i++ {
				var cid = spec.Children[i]
				comp := node.(IComposite)
				comp.AddChild(nodes[cid])
			}
		} else if node.GetCategory() == b3.DECORATOR && len(spec.Child) > 0 {
			dec := node.(IDecorator)
			dec.SetChild(nodes[spec.Child])
		}
	}

	bt.root = nodes[data.Root]
}

func (bt *BehaviorTree) dump() *config.BTTreeCfg {
	return bt.dumpInfo
}

func (bt *BehaviorTree) Tick(target interface{}, blackboard *Blackboard) b3.Status {
	if blackboard == nil {
		panic("The blackboard parameter is obligatory and must be an instance of b3.Blackboard")
	}

	/* CREATE A TICK OBJECT */
	var tick = NewTick()
	tick.debug = bt.debug
	tick.target = target
	tick.Blackboard = blackboard
	tick.tree = bt

	/* TICK NODE */
	var state = bt.root._execute(tick)

	/* CLOSE NODES FROM LAST TICK, IF NEEDED */
	var lastOpenNodes = blackboard._getTreeData(bt.id).OpenNodes
	var currOpenNodes []IBaseNode
	currOpenNodes = append(currOpenNodes, tick._openNodes...)

	// does not close if it is still open in bt tick
	var start = 0
	for i := 0; i < b3.MinInt(len(lastOpenNodes), len(currOpenNodes)); i++ {
		start = i + 1
		if lastOpenNodes[i] != currOpenNodes[i] {
			break
		}
	}

	// close the nodes
	for i := len(lastOpenNodes) - 1; i >= start; i-- {
		lastOpenNodes[i]._close(tick)
	}

	/* POPULATE BLACKBOARD */
	blackboard._getTreeData(bt.id).OpenNodes = currOpenNodes
	blackboard.SetTree("nodeCount", tick._nodeCount, bt.id)

	return state
}

func (bt *BehaviorTree) Print() {
	printNode(bt.root, 0)
}

func printNode(root IBaseNode, blk int) {

	//fmt.Println("new node:", root.Name, " children:", len(root.Children), " child:", root.Child)
	for i := 0; i < blk; i++ {
		fmt.Print(" ") //缩进
	}

	//fmt.Println("|—<", root.Name, ">") //打印"|—<id>"形式
	fmt.Print("|—", root.GetTitle())

	if root.GetCategory() == b3.DECORATOR {
		dec := root.(IDecorator)
		if dec.GetChild() != nil {
			//fmt.Print("=>")
			printNode(dec.GetChild(), blk+3)
		}
	}

	fmt.Println("")
	if root.GetCategory() == b3.COMPOSITE {
		comp := root.(IComposite)
		if comp.GetChildCount() > 0 {
			for i := 0; i < comp.GetChildCount(); i++ {
				printNode(comp.GetChild(i), blk+3)
			}
		}
	}

}
