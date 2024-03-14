package core

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"

	"github.com/dfklegend/cell2/utils/golua"
	"mlrs/b3"
	. "mlrs/b3/config"
)

type IBaseWrapper interface {
	_execute(tick *Tick) b3.Status
	_enter(tick *Tick)
	_open(tick *Tick)
	_tick(tick *Tick) b3.Status
	_close(tick *Tick)
	_exit(tick *Tick)
}
type IBaseNode interface {
	IBaseWrapper

	Ctor()
	Initialize(params *BTNodeCfg)
	GetCategory() string
	Execute(tick *Tick) b3.Status
	GetName() string
	GetTitle() string
	SetBaseNodeWorker(worker IBaseWorker)
	GetBaseNodeWorker() IBaseWorker
}

type BaseNode struct {
	IBaseWorker

	id string

	name string

	category string

	title string

	description string

	parameters map[string]interface{}

	properties map[string]interface{}
}

func (node *BaseNode) Ctor() {

}

func (node *BaseNode) SetName(name string) {
	node.name = name
}
func (node *BaseNode) SetTitle(title string) {
	node.title = title
}

func (node *BaseNode) SetBaseNodeWorker(worker IBaseWorker) {
	node.IBaseWorker = worker
}

func (node *BaseNode) GetBaseNodeWorker() IBaseWorker {
	return node.IBaseWorker
}

func (node *BaseNode) Initialize(params *BTNodeCfg) {
	//node.id = b3.CreateUUID()
	//node.title       = node.title || node.name
	node.description = ""
	node.parameters = make(map[string]interface{})
	node.properties = make(map[string]interface{})

	node.id = params.Id //|| node.id;
	node.name = params.Name
	node.title = params.Title             //|| node.title;
	node.description = params.Description // || node.description;
	node.properties = params.Properties   //|| node.properties;

}

func (node *BaseNode) GetCategory() string {
	return node.category
}

func (node *BaseNode) GetID() string {
	return node.id
}

func (node *BaseNode) GetName() string {
	return node.name
}
func (node *BaseNode) GetTitle() string {
	//fmt.Println("GetTitle ", node.title)
	return node.title
}

func (node *BaseNode) _execute(tick *Tick) b3.Status {
	//fmt.Println("_execute :", node.title)
	// ENTER
	node._enter(tick)

	// OPEN
	if !tick.Blackboard.GetBool("isOpen", tick.tree.id, node.id) {
		node._open(tick)
	}

	// TICK
	var status = node._tick(tick)

	// CLOSE
	if status != b3.RUNNING {
		node._close(tick)
	}

	// EXIT
	node._exit(tick)

	return status
}
func (node *BaseNode) Execute(tick *Tick) b3.Status {
	return node._execute(tick)
}

func (node *BaseNode) _enter(tick *Tick) {
	tick._enterNode(node)
	node.OnEnter(tick)
}

func (node *BaseNode) _open(tick *Tick) {
	//fmt.Println("_open :", node.title)
	tick._openNode(node)
	tick.Blackboard.Set("isOpen", true, tick.tree.id, node.id)

	mem := tick.Blackboard.GetMem("lua")
	if mem != nil {
		luaEngine := mem.(*golua.LuaEngine)

		luaNode := fmt.Sprintf("node/%s.lua", node.name)
		if golua.IsCompiled(luaNode) {
			_, err := luaEngine.DoLuaMethodWithResult(luaNode, "OnOpen", node, tick)
			if err != nil {
				panic(err)
			}
		} else {
			node.OnOpen(tick)
		}
	} else {
		node.OnOpen(tick)
	}
}

func (node *BaseNode) _tick(tick *Tick) b3.Status {
	//fmt.Println("_tick :", node.title)
	tick._tickNode(node)
	var result = b3.FAILURE

	mem := tick.Blackboard.GetMem("lua")
	if mem != nil {
		luaEngine := mem.(*golua.LuaEngine)

		luaNode := fmt.Sprintf("node/%s.lua", node.name)
		if golua.IsCompiled(luaNode) {
			ret, err := luaEngine.DoLuaMethodWithResult(luaNode, "OnTick", node, tick)
			if err != nil {
				panic(err)
			}

			if ret == nil {
				return node.OnTick(tick)
			}

			if ret.Type() != lua.LTString {
				panic(fmt.Sprintf("Node = %s,Method = OnTick, 错误返回值 %s", node.name, ret.String()))
			}

			switch ret.String() {
			case "success":
				result = b3.SUCCESS
			case "failure", "fail":
				result = b3.FAILURE
			case "error":
				result = b3.ERROR
			case "running":
				result = b3.RUNNING
			default:
				result = b3.FAILURE
			}
		} else {
			result = node.OnTick(tick)
		}
	} else {
		result = node.OnTick(tick)
	}

	return result
}

func (node *BaseNode) _close(tick *Tick) {
	tick._closeNode(node)
	tick.Blackboard.Set("isOpen", false, tick.tree.id, node.id)
	node.OnClose(tick)
}

func (node *BaseNode) _exit(tick *Tick) {
	tick._exitNode(node)
	node.OnExit(tick)
}
