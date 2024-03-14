package app

import (
	"errors"
	"fmt"

	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/node/builtin/channel"
	"github.com/dfklegend/cell2/node/builtin/msgs"
	"github.com/dfklegend/cell2/node/service"
	l "github.com/dfklegend/cell2/utils/logger"
)

var (
	ErrorNoService = errors.New("can not find service")
)

// Request
// 发起一次服务请求
// 自动根据路由规则来进行路由，选择目标服务
// route: serviceType.registry.method
func Request(ns *service.NodeService, route string, routeParam any, msg any, cbFunc func(error, any)) {
	serviceType, apiType, method := SplitClientRoute(route)

	pid := RoutePID(serviceType, routeParam)
	if pid == nil {
		l.Log.Errorf("route failed service: route: %v param: %v", route, routeParam)
		apientry.CheckInvokeCBFunc(cbFunc, ErrorNoService, nil)
		return
	}
	ns.RequestEx(pid, fmt.Sprintf("%v.%v", apiType, method), msg, cbFunc)
}

// Notify
// 发起一次通知
func Notify(sender *service.NodeService, route string, routeParam any, msg any) {
	serviceType, apiType, method := SplitClientRoute(route)

	pid := RoutePID(serviceType, routeParam)
	if pid == nil {
		l.Log.Errorf("route failed service: route: %v param: %v", route, routeParam)
		return
	}
	sender.NotifyEx(pid, fmt.Sprintf("%v.%v", apiType, method), msg)
}

// QuerySession
// 向目标front查询session数据
func QuerySession(ns *service.NodeService, frontId string, sessionId uint32, cb func(error, any)) {
	pid := GetServicePID(frontId)
	if pid == nil {
		l.Log.Errorf("query session can not find service: %v", frontId)
		return
	}

	ns.RequestEx(pid, "sys.querysession", &msgs.QuerySession{
		SessionId: sessionId,
	}, cb)
}

// Kick
// 要求目标front踢掉某个连接
func Kick(ns *service.NodeService, frontId string, sessionId uint32, cb func(error, any)) {
	pid := GetServicePID(frontId)
	if pid == nil {
		l.Log.Errorf("Kick can not find service: %v", frontId)
		return
	}

	ns.RequestEx(pid, "sys.kick", &msgs.Kick{
		SessionId: sessionId,
	}, cb)
}

// PushMessageById
// 向一个客户端推送消息
func PushMessageById(ns *service.NodeService, serverId string, sessionId uint32, route string, msg any) {
	channel.GetPushImpl().PushMessageById(ns, serverId, sessionId, route, msg)
}

// PushMessageByIds
// 向多个客户端推送消息
func PushMessageByIds(ns *service.NodeService, serverId string, ids []uint32, route string, msg any) {
	channel.GetPushImpl().PushMessageByIds(ns, serverId, ids, route, msg)
}
