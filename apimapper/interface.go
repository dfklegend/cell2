package apimapper

import (
	"reflect"
)

//	IContext 调用rpc的入口参数
type IContext interface {
	Reserve()
	Handle()
}

type DummyContext struct {
}

func (d *DummyContext) Reserve() {}
func (d *DummyContext) Handle()  {}

// 	API格式器
// 	提供分析API，API调用抽象
// 	前端接口和rpc
// 	receiver, context, inArg, cb
// 	rpc接口没有session
// 	handler接口有session
type IAPIFormatter interface {
	//	IsValidMethod 是否符合接口需求
	IsValidMethod(reflect.Method) bool
}

var (
	TypeOfContext = reflect.TypeOf((*IContext)(nil)).Elem()
	TypeOfError   = reflect.TypeOf((*error)(nil)).Elem()
)

// 	----------------
type IAPIEntry interface {
	Desc() string
}

// 	----------------
type BaseAPIEntry struct {
}

func (b *BaseAPIEntry) Desc() string {
	return "BaseAPIEntry"
}

type APIEntry struct {
	BaseAPIEntry
}

// 	----------------
