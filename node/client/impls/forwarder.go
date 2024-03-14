package impls

import (
	"errors"

	as "github.com/dfklegend/cell2/actorex/service"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/builtin/msgs"
	cs "github.com/dfklegend/cell2/node/client/session"
	"github.com/dfklegend/cell2/node/service"
	l "github.com/dfklegend/cell2/utils/logger"
)

type ForwarderComponent struct {
	*service.BaseComponent
}

func NewForwarder() *ForwarderComponent {
	return &ForwarderComponent{
		BaseComponent: service.NewBaseComponent(),
	}
}

func (f *ForwarderComponent) OnAdd() {
	//l.Log.Infof("handler on add")
	//f.GetNodeService().RequestEx("")
}

func (f *ForwarderComponent) Forward(fs *cs.FrontSession, msg *msgs.ClientMsg) {
	// . 路由找到目标service
	// . request(sys.forward)
	route := msg.Route
	serviceType, _, _ := SplitClientRoute(route)

	// 替换为路由模块
	pid := app.RoutePID(serviceType, fs)
	if pid == nil {
		l.Log.Errorf("forwarder error, can not find target service, route: %v, check route params(is pushsession succ?)", route)
		l.Log.Errorf("fs: %v", fs.ToJson())
		return
	}

	session := fs.Session

	msg.ID = fs.GetID()
	msg.FrontId = f.GetNodeService().Name

	// make session data
	//msg.SessionData = []byte(fs.ToJson())

	isNotify := msg.ClientReqId == as.NotifyReqID
	if isNotify {
		f.GetNodeService().NotifyEx(pid, "sys.notify", msg)
		return
	}

	// todo: forward调用的错误和返回的错误如何区分开
	// err 应该代表forward本身产生的错误
	// Response.Error则是接口返回的错误，应该传递回去的
	f.GetNodeService().RequestEx(pid, "sys.call", msg, func(err error, ret any) {
		if err != nil {
			l.Log.Errorf("%v forward error: %v", route, err)
			session.ResponseMID(uint(msg.ClientReqId), nil, err)
			return
		}

		res, _ := ret.(*msgs.Response)
		if res == nil {
			l.Log.Errorf("ret is not Response")
			return
		}

		if res.SessionId != session.GetId() || res.ClientReqId != msg.ClientReqId {
			l.Log.Errorf("missmatch res")
			return
		}

		if res.Error != "" {
			l.Log.Errorf("%v forward result error str: %v", route, res.Error)
			session.ResponseMID(uint(msg.ClientReqId), nil, errors.New(res.Error))
		} else {
			session.ResponseMID(uint(msg.ClientReqId), res.Data, nil)
		}

	})
}
