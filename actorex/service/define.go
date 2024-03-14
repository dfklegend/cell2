package service

import (
	"github.com/asynkron/protoactor-go/actor"

	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	"github.com/dfklegend/cell2/utils/runservice"
)

const (
	CodeSucc int32 = iota // 成功

	CodeErrString = 999 // 一般性错误，错误字符串

	// 	reserved range(1-1000)
	CodeUserBegin = 1000 // 用户code
)

const (
	// NotifyReqID reqId 为0代表通知消息
	NotifyReqID = 0
)

// 	回调函数
//	Error, 返回值对象
type ResCBFunc func(error, interface{})

type IRequestReceiver interface {
	ReceiveRequest(ctx actor.Context, request *messages.ServiceRequest, rawMsg interface{})
}

type IService interface {
	SetExtProps(props *ExtProps)
	OnCreate()
	SetRunService(rservice *runservice.StandardRunService)
	GetRunService() *runservice.StandardRunService
	SetAPIDispatcher(disp IAPIDispatcher)
}

type IAPIDispatcher interface {
	Dispatch(ctx actor.Context, srv *Service, route string, request *messages.ServiceRequest) bool
}
