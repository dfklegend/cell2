package core

import (
	"fmt"
	"reflect"
)

type TreeData struct {
	NodeMemory     *Memory
	OpenNodes      []IBaseNode
	TraversalDepth int
	TraversalCycle int
}

func NewTreeData() *TreeData {
	return &TreeData{NewMemory(), make([]IBaseNode, 0), 0, 0}
}

// Memory -------------------------------------------------
type Memory struct {
	_memory map[string]interface{}
}

func NewMemory() *Memory {
	return &Memory{make(map[string]interface{})}
}

func (b *Memory) Get(key string) interface{} {
	return b._memory[key]
}
func (b *Memory) Set(key string, val interface{}) {
	b._memory[key] = val
}
func (b *Memory) Remove(key string) {
	delete(b._memory, key)
}

// TreeMemory -------------------------------------------------
type TreeMemory struct {
	*Memory
	_treeData   *TreeData
	_nodeMemory map[string]*Memory
}

func NewTreeMemory() *TreeMemory {
	return &TreeMemory{NewMemory(), NewTreeData(), make(map[string]*Memory)}
}

// Blackboard -------------------------------------------------
type Blackboard struct {
	_baseMemory *Memory
	_treeMemory map[string]*TreeMemory
}

func NewBlackboard() *Blackboard {
	p := &Blackboard{}
	p.Initialize()
	return p
}

func (b *Blackboard) Initialize() {
	b._baseMemory = NewMemory()
	b._treeMemory = make(map[string]*TreeMemory)
}

func (b *Blackboard) _getTreeMemory(treeScope string) *TreeMemory {
	if _, ok := b._treeMemory[treeScope]; !ok {
		b._treeMemory[treeScope] = NewTreeMemory()
	}
	return b._treeMemory[treeScope]
}

func (b *Blackboard) _getNodeMemory(treeMemory *TreeMemory, nodeScope string) *Memory {
	memory := treeMemory._nodeMemory
	if _, ok := memory[nodeScope]; !ok {
		memory[nodeScope] = NewMemory()
	}

	return memory[nodeScope]
}

func (b *Blackboard) _getMemory(treeScope, nodeScope string) *Memory {
	var memory = b._baseMemory

	if len(treeScope) > 0 {
		treeMem := b._getTreeMemory(treeScope)
		memory = treeMem.Memory
		if len(nodeScope) > 0 {
			memory = b._getNodeMemory(treeMem, nodeScope)
		}
	}

	return memory
}

func (b *Blackboard) Set(key string, value interface{}, treeScope, nodeScope string) {
	var memory = b._getMemory(treeScope, nodeScope)
	memory.Set(key, value)
}

func (b *Blackboard) SetMem(key string, value interface{}) {
	var memory = b._getMemory("", "")
	memory.Set(key, value)
}

func (b *Blackboard) Remove(key string) {
	var memory = b._getMemory("", "")
	memory.Remove(key)
}
func (b *Blackboard) SetTree(key string, value interface{}, treeScope string) {
	var memory = b._getMemory(treeScope, "")
	memory.Set(key, value)
}
func (b *Blackboard) _getTreeData(treeScope string) *TreeData {
	treeMem := b._getTreeMemory(treeScope)
	return treeMem._treeData
}

func (b *Blackboard) Get(key, treeScope, nodeScope string) interface{} {
	memory := b._getMemory(treeScope, nodeScope)
	return memory.Get(key)
}
func (b *Blackboard) GetMem(key string) interface{} {
	memory := b._getMemory("", "")
	return memory.Get(key)
}
func (b *Blackboard) GetFloat64(key, treeScope, nodeScope string) float64 {
	v := b.Get(key, treeScope, nodeScope)
	if v == nil {
		return 0
	}
	return v.(float64)
}
func (b *Blackboard) GetBool(key, treeScope, nodeScope string) bool {
	v := b.Get(key, treeScope, nodeScope)
	if v == nil {
		return false
	}
	return v.(bool)
}
func (b *Blackboard) GetInt(key, treeScope, nodeScope string) int {
	v := b.Get(key, treeScope, nodeScope)
	if v == nil {
		return 0
	}
	return v.(int)
}
func (b *Blackboard) GetInt64(key, treeScope, nodeScope string) int64 {
	v := b.Get(key, treeScope, nodeScope)
	if v == nil {
		return 0
	}
	return v.(int64)
}
func (b *Blackboard) GetUInt64(key, treeScope, nodeScope string) uint64 {
	v := b.Get(key, treeScope, nodeScope)
	if v == nil {
		return 0
	}
	return v.(uint64)
}

func (b *Blackboard) GetInt64Safe(key, treeScope, nodeScope string) int64 {
	v := b.Get(key, treeScope, nodeScope)
	if v == nil {
		return 0
	}
	return ReadNumberToInt64(v)
}
func (b *Blackboard) GetUInt64Safe(key, treeScope, nodeScope string) uint64 {
	v := b.Get(key, treeScope, nodeScope)
	if v == nil {
		return 0
	}
	return ReadNumberToUInt64(v)
}

func (b *Blackboard) GetInt32(key, treeScope, nodeScope string) int32 {
	v := b.Get(key, treeScope, nodeScope)
	if v == nil {
		return 0
	}
	return v.(int32)
}

func ReadNumberToInt64(v interface{}) int64 {
	var ret int64
	switch tvalue := v.(type) {
	case uint64:
		ret = int64(tvalue)
	default:
		panic(fmt.Sprintf("错误的类型转成Int64 %v:%+v", reflect.TypeOf(v), v))
	}

	return ret
}

func ReadNumberToUInt64(v interface{}) uint64 {
	var ret uint64
	switch tvalue := v.(type) {
	case int64:
		ret = uint64(tvalue)
	default:
		panic(fmt.Sprintf("错误的类型转成UInt64 %v:%+v", reflect.TypeOf(v), v))
	}
	return ret
}
