package chat

import (
	"strings"
	"time"

	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	nodeapi "github.com/dfklegend/cell2/apimapper/registry"
	"github.com/dfklegend/cell2/node/client/impls"
	l "github.com/dfklegend/cell2/utils/logger"
	mymsg "server/messages"
	"server/messages/clientmsg"
)

func init() {
	nodeapi.Registry.AddCollection("chat.handler").
		Register(&Handler{}, apientry.WithGroupName("chat"), apientry.WithNameFunc(strings.ToLower))
}

var (
	unequalOccured = 0
)

type Handler struct {
	api.APIEntry
}

func (e *Handler) Hello(ctx *impls.HandlerContext, msg *clientmsg.Hello, cbFunc apientry.HandlerCBFunc) error {
	s := ctx.ActorContext.Actor().(*Service)

	//l.Log.Infof("hello in %v\n", s.Name)
	//l.Log.Infof("msg: %+v\n", msg)
	//l.Log.Infof("sessionData: %v", ctx.Session.ToJson())

	//bs := ctx.Session

	//s.GetRunService().GetTimerMgr().After(time.Millisecond, func(args ...any) {
	//	if msg.Log > 0 {
	//		l.L.Infof("got request: %v", ctx.ClientReqId)
	//	}
	//	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
	//		Result: msg.Msg,
	//		Code:   int32(msg.Number),
	//	})
	//})

	clientReqId := ctx.ClientReqId
	if msg.Log > 0 {
		l.L.Infof("got request: %v", clientReqId)
	}

	s.GetRunService().GetTimerMgr().After(time.Millisecond, func(args ...any) {
		// 现在context每次调用都有专属的，下面错误不会发生
		if clientReqId != ctx.ClientReqId {
			unequalOccured++
			if unequalOccured%10000 == 0 {
				l.L.Infof("unequalOccured: %v", unequalOccured)
			}
			panic("unequalOccured")
		}
	})

	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{
		Result: msg.Msg,
		Code:   int32(msg.Number),
	})
	return nil
}
