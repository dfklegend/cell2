package service

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"

	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	l "github.com/dfklegend/cell2/utils/logger"
)

// DirectSendNotify 直接向目标投递一个Notify
// Notify不需要回执，所以存在直接notify的需求
func DirectSendNotify(ctx *actor.RootContext, pid *actor.PID, route string, msg interface{}) {
	bytes, typeName, err := remote.Serialize(msg, DefaultSerializeId)
	if err != nil {
		l.Log.Errorf("msg serialize failed: %v", err)
		return
	}

	req := &messages.ServiceRequest{
		Sender: nil,
		ReqId:  NotifyReqID,
		Route:  route,
		Type:   typeName,
		Body:   bytes,
	}

	ctx.Send(pid, req)
}
