package gate

import (
	"errors"
	"fmt"
	"strings"

	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/apimapper/registry"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/client/impls"
	"github.com/dfklegend/cell2/utils/waterfall"

	"mmo/common/define"
	mymsg "mmo/messages"
	"mmo/messages/cproto"
)

func init() {
	registry.Registry.AddCollection("gate.handler").
		Register(&Handler{}, apientry.WithName("gate"), apientry.WithNameFunc(strings.ToLower))
}

func Visit() {
}

type Handler struct {
	api.APIEntry
}

// QueryGate
// 负载均衡分配gate
func (e *Handler) QueryGate(ctx *impls.HandlerContext, msg *cproto.QueryGateReq, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	ns := s.GetNodeService()
	log := ns.GetLogger()

	log.Infof("queryGate in %v", s.Name)

	// TODO: 此逻辑依赖于配置中本地有所有gate端口配置
	// 改造成,所有gate向gate-1同步自己监听端口，更合理
	connectorId := s.allocGate()
	cfg := app.Node.GetServiceCfg(connectorId)
	if cfg == nil {
		apientry.CheckInvokeCBFunc(cbFunc, errors.New("can not find connector"), nil)
		return
	}

	ip := ""
	port1 := ""
	port2 := ""

	port1 = ""
	subs := strings.Split(cfg.WSClientAddress, ":")
	if len(subs) == 2 {
		ip = subs[0]
		port1 = subs[1]
	}

	subs = strings.Split(cfg.ClientAddress, ":")
	if len(subs) == 2 {
		ip = subs[0]
		port2 = subs[1]
	}

	apientry.CheckInvokeCBFunc(cbFunc, nil,
		&cproto.QueryGateAck{
			Code: 0,
			IP:   ip,
			// wsPort, tcpPort
			Port: fmt.Sprintf("%v,%v", port1, port2),
		})
}

//	选择一个logic，向其发送进入请求
func (e *Handler) Login(ctx *impls.HandlerContext, msg *cproto.LoginReq, cbFunc apientry.HandlerCBFunc) {
	s := ctx.ActorContext.Actor().(*Service)
	ns := s.GetNodeService()
	fs := ctx.Session

	log := ns.GetLogger()

	log.Infof("login in %v", s.Name)
	log.Infof("%+v", msg)

	scheduler := ns.GetRunService().GetScheduler()

	// 公共变量
	var uid int64
	var logicId string
	centerLogined := false

	errCode := define.ErrFaild

	waterfall.NewBuilder(scheduler).Next(func(callback waterfall.Callback, args ...interface{}) {
		// . auth
		app.Request(ns, "db.dbremote.auth", nil, &mymsg.DBAuth{
			Username: msg.Username,
			Password: msg.Password,
		}, func(err error, raw interface{}) {
			ack := raw.(*mymsg.DBAuthAck)
			uid = ack.UId
			// 说明认证失败
			if uid == 0 {
				log.Infof("deremote.auth failed")
				errCode = define.ErrAuthFailed
				callback(true)
				return
			}
			log.Infof("got uid: %v", ack.UId)
			callback(false)
		})

	}).Next(func(callback waterfall.Callback, args ...interface{}) {
		// . 去center reqLogin

		app.Request(ns, "center.centerremote.reqlogin", nil, &mymsg.CenterReqLogin{
			ServerId: s.Name,
			NetId:    ctx.Session.GetNetId(),
			UId:      uid,
			KickPrev: true,
		}, func(err error, raw interface{}) {
			if err != nil {
				log.Errorf("CenterReqLogin error: %v!", err)
				callback(true)
				return
			}
			ack := raw.(*mymsg.CenterReqLoginAck)
			if ack.Code != int32(define.Succ) {
				log.Errorf("CenterReqLogin faild!")
				callback(true)
				return
			}

			log.Infof("CenterReqLogin succ!")
			centerLogined = true
			// 此时事务锁locked(logining 或者 reonline)

			if ack.IsReconnect {
				doReconnect(ctx, ns, uid, ack.LogicId, cbFunc)
				return
			}

			callback(false)
		})

	}).Next(func(callback waterfall.Callback, args ...interface{}) {
		doTestDelay(ns, callback)
	}).Next(func(callback waterfall.Callback, args ...interface{}) {

		// 分配logic
		logicId = app.RandGetWorkServiceName("logic")
		if logicId == "" {
			log.Errorf("find no logic_old")
			callback(true)
			return
		}

		app.Request(ns, "logic.logicremote.enter", logicId, &mymsg.LogicLoadPlayer{
			ServerId: ns.Name,
			NetId:    ctx.Session.GetNetId(),
			UId:      uid,
		}, func(err error, raw interface{}) {
			if err != nil {
				log.Errorf("load player err:%v", err)
				callback(true)
				return
			}
			ack := raw.(*mymsg.NormalAck)
			if ack.Code == int32(define.Succ) {
				// succ
				log.Infof("loadplayer succ")
				callback(false)
				return
			}

			callback(true)
		})
	}).Next(func(callback waterfall.Callback, args ...interface{}) {
		// 判断一下是否掉线了
		if fs.IsClosed() {
			log.Infof("session closed already!")
			callback(true)
			return
		}

		fs.Bind(fmt.Sprintf("%v", uid))
		fs.Set("logic", logicId)
		// 登录成功了，设置一下uid,便于掉线通知
		fs.Set("uid", uid)
		fs.PushSession(nil)

		callback(false)

	}).Final(func(err bool, args ...interface{}) {

		log.Infof("final isError:%v", err)

		if err {
			// 出错的处理
			if centerLogined {
				sendSessionClose(ns, uid)
			}

			if !fs.IsClosed() {
				apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.LoginAck{
					Code: int32(errCode),
				})
			}
			return
		}

		apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.LoginAck{
			UId: uid,
		})

	}).Do()
}
