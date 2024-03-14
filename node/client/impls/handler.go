package impls

import (
	"errors"
	"fmt"

	as "github.com/dfklegend/cell2/actorex/service"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/node/builtin/msgs"
	"github.com/dfklegend/cell2/node/client"
	"github.com/dfklegend/cell2/node/client/impls/config"
	cs "github.com/dfklegend/cell2/node/client/session"
	"github.com/dfklegend/cell2/node/service"
	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/sche"
	"github.com/dfklegend/cell2/utils/serialize"
)

// HandlerComponent
// 处理来自客户端的消息
// 映射到脚本
// 本服务
// clients -> handlerComponent -> handlers
// 远程服务
// clients -> handlerComponent -> forwarderComponent -> (target)handlerComponent -> handlers
type HandlerComponent struct {
	*service.BaseComponent
	api *apientry.APICollection

	scheduler *sche.Sche
	forwarder *ForwarderComponent

	onCloseCBs map[uint32]client.CBSessionOnClose
}

func NewHandler(api *apientry.APICollection) *HandlerComponent {
	return &HandlerComponent{
		BaseComponent: service.NewBaseComponent(),
		api:           api,
		onCloseCBs:    make(map[uint32]client.CBSessionOnClose),
	}
}

func (h *HandlerComponent) SetForwarder(f *ForwarderComponent) {
	h.forwarder = f
}

func (h *HandlerComponent) OnAdd() {
	//l.Log.Infof("handler on add")
	h.scheduler = h.GetNodeService().GetRunService().GetScheduler()
}

func (h *HandlerComponent) OnRemove() {
}

func (h *HandlerComponent) Process(fs *cs.FrontSession, msg *msgs.ClientMsg) {
	route := msg.Route
	serviceType, apiType, method := SplitClientRoute(route)

	ns := h.GetNodeService()
	if serviceType != ns.GetServiceType() {
		//l.Log.Infof("forward to other service")
		h.forwarder.Forward(fs, msg)
		return
	}

	session := fs.Session
	apiRoute := fmt.Sprintf("%v.%v", apiType, method)

	serializer := config.GetConfig().Serializer
	ctx := h.allocContext()
	ctx.Update(ns.Context, fs, msg.ClientReqId)

	// 也需要注意session可能在某个异步操作后，失效

	h.tryCallCol(ctx, msg, h.api, apiRoute, serializer, func(err error, ret any) {
		if err != nil {
			session.ResponseMID(uint(msg.ClientReqId), nil, err)
			return
		}

		// return
		retData, serializeErr := serializer.Marshal(ret)
		if serializeErr != nil {
			session.ResponseMID(uint(msg.ClientReqId), nil, serializeErr)
			return
		}
		session.ResponseMID(uint(msg.ClientReqId), retData, nil)
	})
}

// 	return
//		process by this collection
//		error of call
func (h *HandlerComponent) tryCallCol(
	ctx *HandlerContext, msg *msgs.ClientMsg,
	col *apientry.APICollection, route string,
	serializer serialize.Serializer, cb func(error, any)) bool {
	if !col.HasMethod(route) {
		l.Log.Errorf("no method:%v", route)
		cb(errors.New("no method"), nil)
		return false
	}

	if msg.ClientReqId == as.NotifyReqID {
		// notify
		apientry.CallWithSerialize(col, ctx, route, msg.Data, nil, serializer)
		return true
	}
	apientry.CallWithSerialize(col, ctx, route, msg.Data, func(err error, ret interface{}) {
		if err != nil {
			l.Log.Errorf("call error:%v", err)
			cb(err, nil)
			return
		}

		cb(nil, ret)
		// 返回值
	}, serializer)
	return true
}

func (h *HandlerComponent) allocContext() *HandlerContext {
	return NewHandlerContext()
}

func (h *HandlerComponent) allocBackSession() *cs.BackSession {
	return cs.NewBS()
}

// ProcessForwardMsg
// 处理被转发过来的消息
func (h *HandlerComponent) ProcessForwardMsg(msg *msgs.ClientMsg, cbFunc apientry.HandlerCBFunc) {
	route := msg.Route
	serviceType, apiType, method := SplitClientRoute(route)

	ns := h.GetNodeService()
	if serviceType != ns.GetServiceType() {
		l.Log.Errorf("got msg: %v at wrong service: %v", route, ns.Name)
		return
	}

	apiRoute := fmt.Sprintf("%v.%v", apiType, method)
	serializer := config.GetConfig().Serializer

	bs := h.allocBackSession()
	cs.InitBackSession(bs, ns, msg.FrontId, msg.SessionId, msg.ID)
	ctx := h.allocContext()
	ctx.Update(ns.Context, bs, msg.ClientReqId)

	h.tryCallCol(ctx, msg, h.api, apiRoute, serializer, func(err error, ret any) {
		if err != nil {
			// 接口错误
			apientry.CheckInvokeCBFunc(cbFunc, nil, &msgs.Response{
				SessionId:   msg.SessionId,
				ClientReqId: msg.ClientReqId,
				Error:       err.Error(),
			})
			return
		}

		data, err := serializer.Marshal(ret)
		apientry.CheckInvokeCBFunc(cbFunc, nil, &msgs.Response{
			SessionId:   msg.SessionId,
			ClientReqId: msg.ClientReqId,
			Data:        data,
		})
	})
}

func (h *HandlerComponent) OnSessionAdd(fs *cs.FrontSession) {

}

func (h *HandlerComponent) OnSessionRemove(fs *cs.FrontSession) {
	onClose := h.onCloseCBs[fs.GetNetId()]
	if onClose == nil {
		return
	}
	onClose(h.GetNodeService(), fs)
	delete(h.onCloseCBs, fs.GetNetId())
}

func (h *HandlerComponent) AddOnSessionClose(netId uint32, onClose client.CBSessionOnClose) {
	h.onCloseCBs[netId] = onClose
}
