package chat

import (
	"strings"

	mymsg "chat2/messages"
	"chat2/messages/clientmsg"
	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	nodeapi "github.com/dfklegend/cell2/apimapper/registry"
	"github.com/dfklegend/cell2/node/client/impls"
	"github.com/dfklegend/cell2/node/client/session"
	l "github.com/dfklegend/cell2/utils/logger"
)

func init() {
	nodeapi.Registry.AddCollection("chat.handler").
		Register(&Handler{}, apientry.WithGroupName("chat"), apientry.WithNameFunc(strings.ToLower))
}

var (
	testIndex = 0
)

type Handler struct {
	api.APIEntry
}

func (e *Handler) SendChat(ctx *impls.HandlerContext, msg *clientmsg.ChatMsg, cbFunc apientry.HandlerCBFunc) error {
	s := ctx.ActorContext.Actor().(*Service)

	l.Log.Infof("SendChat in %v\n", s.Name)
	l.Log.Infof("msg: %+v\n", msg)
	l.Log.Infof("before query sessionData: %v", ctx.Session.ToJson())

	bs := ctx.Session.(*session.BackSession)
	bs.QuerySession(func(err error) {
		l.Log.Infof("query session ret")
		l.Log.Infof("sessionData: %v", bs.ToJson())
	})

	testIndex++
	bs.Set("testindex", testIndex)
	bs.PushSession(nil)

	s.cs.OnChat(bs.GetID(), msg.Content)
	apientry.CheckInvokeCBFunc(cbFunc, nil, &mymsg.NormalAck{})

	return nil
}
