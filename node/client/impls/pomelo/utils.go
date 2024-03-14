package pomelo

import (
	"reflect"
	"time"

	as "github.com/dfklegend/cell2/actorex/service"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/pomelonet/server/acceptor"
	"github.com/dfklegend/cell2/pomelonet/server/session"
	"github.com/dfklegend/cell2/utils/logger"
)

func StartAcceptor(a acceptor.Acceptor, cfg *session.SessionConfig) {
	// 监听
	go func() {
		a.ListenAndServe()
	}()

	// 获取新连接
	go func() {
		for conn := range a.GetConnChan() {
			logger.Log.Debugf("new conn come:%v", conn)
			// 新连接建立
			s := session.NewClientSession(conn, cfg)
			// add to
			s.Handle()
			logger.Log.Debugf("%v", conn)
		}
	}()

	go func() {
		time.Sleep(time.Second)
		logger.Log.Infof("listening with acceptor %s on addr %s", reflect.TypeOf(a), a.GetAddr())
	}()
}

func TryCreateAcceptors(s as.IService, name string) {
	ns, _ := s.(service.INodeServiceOwner)
	ServiceCreateAcceptors(ns.GetNodeService(), name, app.Node.GetServiceCfg(name))
}
