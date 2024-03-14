package impls

import (
	"github.com/asynkron/protoactor-go/actor"

	cs "github.com/dfklegend/cell2/node/client/session"
)

type HandlerContext struct {
	ActorContext actor.Context
	ClientReqId  uint32
	// FrontSession or BackSession
	// 实际前端服务器 FrontSession
	// 后端服务器 BackSession
	Session cs.IServerSession
}

func NewHandlerContext() *HandlerContext {
	return &HandlerContext{}
}

func (r *HandlerContext) Update(ctx actor.Context, session cs.IServerSession, clientReqId uint32) {
	r.ActorContext = ctx
	r.Session = session
	r.ClientReqId = clientReqId
}

func (r *HandlerContext) SetClientReqId(id uint32) {
	r.ClientReqId = id
}

func (r *HandlerContext) GetFrontSession() *cs.FrontSession {
	return r.Session.(*cs.FrontSession)
}

func (r *HandlerContext) GetBackSession() *cs.BackSession {
	return r.Session.(*cs.BackSession)
}

func (r *HandlerContext) GetSession() cs.IServerSession {
	return r.Session
}

func (r *HandlerContext) reset() {
	r.ActorContext = nil
	r.Session = nil
	r.ClientReqId = 0
}

func (r *HandlerContext) Reserve() {
}

func (r *HandlerContext) Handle() {
}
