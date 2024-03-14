package builtin

import (
	"strings"

	as "github.com/dfklegend/cell2/actorex/service"
	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/apimapper/registry"
	"github.com/dfklegend/cell2/node/builtin/msgs"
	"github.com/dfklegend/cell2/node/client"
	"github.com/dfklegend/cell2/node/client/impls"
	"github.com/dfklegend/cell2/node/service"
	l "github.com/dfklegend/cell2/utils/logger"
)

func init() {
	registry.Registry.AddCollection(service.SystemAPI).
		Register(&Entry{}, apientry.WithGroupName("sys"), apientry.WithNameFunc(strings.ToLower))
}

func Visit() {
}

//	系统接口
//	比如，转发前端消息
type Entry struct {
	api.APIEntry
}

//	处理转发的notify消息
func (e *Entry) Notify(ctx *as.RemoteContext, msg *msgs.ClientMsg) {
	s := ctx.ActorContext.Actor().(service.INodeServiceOwner)
	ns := s.GetNodeService()

	// 调用handler
	handler := ns.GetComponent("handler").(*impls.HandlerComponent)
	if handler == nil {
		l.Log.Errorf("can not find handler")
		return
	}

	if msg.ClientReqId != as.NotifyReqID {
		l.Log.Errorf("bad call sys.notify")
	}

	handler.ProcessForwardMsg(msg, nil)
	return
}

// 	处理转发的Request
func (e *Entry) Call(ctx *as.RemoteContext, msg *msgs.ClientMsg, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(service.INodeServiceOwner)
	ns := s.GetNodeService()

	// 调用handler
	handler := ns.GetComponent("handler").(*impls.HandlerComponent)
	if handler == nil {
		l.Log.Errorf("can not find handler")
		return
	}

	handler.ProcessForwardMsg(msg, cbFunc)
	return
}

func (e *Entry) PushMsg(ctx *as.RemoteContext, msg *msgs.PushMsg, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(service.INodeServiceOwner)
	ns := s.GetNodeService()

	sc := ns.GetComponent("sessions").(*impls.SessionsComponent)

	if sc == nil {
		l.Log.Errorf("can not find sessions")
		return
	}

	sessions := sc.GetSessions()
	sessions.PushMsg(msg)
	apientry.CheckInvokeCBFunc(cbFunc, nil, nil)
	return
}

func (e *Entry) PushSession(ctx *as.RemoteContext, msg *msgs.PushSession, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(service.INodeServiceOwner)
	ns := s.GetNodeService()

	sc := ns.GetComponent("sessions").(*impls.SessionsComponent)
	if sc == nil {
		l.Log.Errorf("can not find sessions")
		return
	}

	sessions := sc.GetSessions()
	sessions.PushSession(msg.SessionId, msg.SessionData)

	apientry.CheckInvokeCBFunc(cbFunc, nil, nil)
	return
}

func (e *Entry) QuerySession(ctx *as.RemoteContext, msg *msgs.QuerySession, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(service.INodeServiceOwner)
	ns := s.GetNodeService()

	sc := ns.GetComponent("sessions").(*impls.SessionsComponent)
	if sc == nil {
		l.Log.Errorf("can not find sessions")
		return
	}

	sessions := sc.GetSessions()
	fs := sessions.GetSession(msg.SessionId)
	if fs == nil {
		apientry.CheckInvokeCBFunc(cbFunc, client.ErrorNoSession, nil)
		return
	}

	ack := &msgs.QuerySessionAck{
		SessionId:   msg.SessionId,
		SessionData: []byte(fs.ToJson()),
	}
	apientry.CheckInvokeCBFunc(cbFunc, nil, ack)
	return
}

func (e *Entry) Kick(ctx *as.RemoteContext, msg *msgs.Kick, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(service.INodeServiceOwner)
	ns := s.GetNodeService()

	sc := ns.GetComponent("sessions").(*impls.SessionsComponent)
	if sc == nil {
		l.Log.Errorf("can not find sessions")
		return
	}

	// 为了让客户端有响应，应该发送一个踢人消息再延迟1s断开
	// 这个属于定制实现，因为协议应该是无关的，由sessions内kickHandler实现

	sessions := sc.GetSessions()
	sessions.Kick(msg.SessionId)

	apientry.CheckInvokeCBFunc(cbFunc, nil, nil)
	return
}
