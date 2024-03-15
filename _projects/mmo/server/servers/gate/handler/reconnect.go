package gate

import (
	"fmt"

	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/client/impls"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/waterfall"

	"mmo/common/define"
	mymsg "mmo/messages"
	"mmo/messages/cproto"
)

// doReconnect
// 断线重连处理
func doReconnect(ctx *impls.HandlerContext, ns *service.NodeService, uid int64, logicId string, cbFunc apientry.HandlerCBFunc) {
	// . 向logic要求ReOnline
	// . bind logicid
	fs := ctx.Session
	log := ns.GetLogger()
	log.Infof("doReconnect: %v %v", uid, logicId)

	waterfall.NewBuilder(ns.GetRunService().GetScheduler()).
		Next(func(callback waterfall.Callback, args ...interface{}) {
			doTestDelay(ns, callback)
		}).Next(func(callback waterfall.Callback, args ...interface{}) {
		app.Request(ns, "logic.logicremote.reenter", logicId, &mymsg.LogicLoadPlayer{
			ServerId: ns.Name,
			NetId:    fs.GetNetId(),
			UId:      uid,
		}, func(err error, raw interface{}) {
			if err != nil {
				log.Errorf("reenter err:%v", err)
				callback(true)
				return
			}
			ack := raw.(*mymsg.NormalAck)
			if ack.Code == int32(define.Succ) {
				// succ
				log.Infof("reenter succ")
				callback(false)
				return
			}

			callback(true)
		})
	}).
		Next(func(callback waterfall.Callback, args ...interface{}) {
			// 判断一下是否掉线了
			if fs.IsClosed() {
				log.Infof("session closed already!")
				callback(true)
				return
			}

			fs.Bind(fmt.Sprintf("%v", uid))
			fs.Set("logic", logicId)
			fs.Set("uid", uid)
			fs.PushSession(nil)

			callback(false)
		}).
		Final(func(err bool, args ...interface{}) {
			log.Infof("final isError:%v", err)

			if err {
				sendSessionClose(ns, uid)

				if !fs.IsClosed() {
					apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.LoginAck{
						Code: int32(define.ErrFaild),
					})
				}
				return
			}

			apientry.CheckInvokeCBFunc(cbFunc, nil, &cproto.LoginAck{
				UId:         uid,
				IsReconnect: true,
			})
		}).
		Do()
}
