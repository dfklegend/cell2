package gate

import (
	"errors"
	"fmt"
	"log"
	"strings"

	mymsg "chat2/messages"
	"chat2/messages/clientmsg"
	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	nodeapi "github.com/dfklegend/cell2/apimapper/registry"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/client/impls"
	cs "github.com/dfklegend/cell2/node/client/session"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/common"
	l "github.com/dfklegend/cell2/utils/logger"
)

func init() {
	nodeapi.Registry.AddCollection("gate.handler").
		Register(&Handler{}, apientry.WithGroupName("gate"), apientry.WithNameFunc(strings.ToLower))
}

func Visit() {
}

type Handler struct {
	api.APIEntry
}

//	分配另外的gate
func (e *Handler) QueryGate(ctx *impls.HandlerContext, msg *clientmsg.QueryGateReq, cbFunc apientry.HandlerCBFunc) error {
	s := ctx.ActorContext.Actor().(*Service)
	log.Printf("queryGate in %v\n", s.Name)

	// 此逻辑依赖于本地由所有gate端口配置
	// 改造成,所有gate向gate-1注册自己，更合理
	connectorId := s.allocGate()
	cfg := app.Node.GetServiceCfg(connectorId)
	if cfg == nil {
		apientry.CheckInvokeCBFunc(cbFunc, errors.New("can not find connector"), nil)
		return nil
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
		&clientmsg.QueryGateAck{
			Code: 0,
			IP:   ip,
			// wsPort, tcpPort
			Port: fmt.Sprintf("%v,%v", port1, port2),
		})

	return nil
}

//	选择一个chat，向其发送进入请求
func (e *Handler) Login(ctx *impls.HandlerContext, msg *clientmsg.LoginReq, cbFunc apientry.HandlerCBFunc) error {
	s := ctx.ActorContext.Actor().(*Service)
	l.Log.Infof("login in %v", s.Name)
	l.Log.Infof("%+v", msg)

	fs := ctx.GetFrontSession()
	impls.AddOnSessionOnClose(s.GetNodeService(), fs.GetNetId(), onSessionClose)

	chatItem := app.RandGetServiceItem("chat")
	if chatItem == nil {
		l.Log.Info("can not find chat")
		return nil
	}
	//pid := chatItem.PID

	uid := fmt.Sprintf("uid-%v-%v", s.Name, fs.GetNetId())
	fs.Bind(uid)
	fs.Set("chatid", chatItem.Name)
	fs.PushSession(nil)

	//s.RequestEx(pid, "chatremote.entry", &mymsg.RoomEntry{
	//	UId:      uid,
	//	Name:     msg.Name,
	//	ServerId: s.Name,
	//	NetId:    ctx.Session.GetNetId(),
	//}, func(err error, ret any) {
	//	code := 0
	//	errStr := ""
	//	if err != nil {
	//		code = 1
	//		errStr = err.Error()
	//	}
	//	e.doLoginRet(cbFunc, code, errStr)
	//})

	app.Request(s.GetNodeService(), "chat.chatremote.entry", fs, &mymsg.RoomEntry{
		UId:      uid,
		Name:     msg.Name,
		ServerId: s.Name,
		NetId:    ctx.Session.GetNetId(),
	}, func(err error, ret any) {
		code := 0
		errStr := ""
		if err != nil {
			code = 1
			errStr = err.Error()
		}
		e.doLoginRet(cbFunc, code, errStr)
	})
	return nil
}

func (e *Handler) doLoginRet(cbFunc apientry.HandlerCBFunc, code int, errStr string) {
	apientry.CheckInvokeCBFunc(cbFunc, nil,
		&clientmsg.NormalAck{
			Code:   code,
			Result: errStr,
		})
}

func onSessionClose(ns *service.NodeService, fs *cs.FrontSession) {
	l.Log.Infof("routine:%v on session close: %v %v",
		common.GetRoutineID(), ns.Name, fs.GetNetId())

	chatId := fs.Get("chatid", "").(string)
	if chatId == "" {
		return
	}

	l.Log.Infof("found chatId: %v", chatId)
	//pid := app.GetServicePID(chatId)
	//if pid == nil {
	//	return
	//}
	//
	//ns.RequestEx(pid, "chatremote.leave", &mymsg.RoomLeave{
	//	UId:      fs.GetID(),
	//	ServerId: ns.Name,
	//	NetId:    fs.GetNetId(),
	//}, func(err error, ret any) {
	//	if err != nil {
	//		l.Log.Infof("chatremote.leave ret, err: %v", err)
	//		return
	//	}
	//	l.Log.Infof("chatremote.leave ret")
	//})

	app.Request(ns, "chat.chatremote.leave", fs, &mymsg.RoomLeave{
		UId:      fs.GetID(),
		ServerId: ns.Name,
		NetId:    fs.GetNetId(),
	}, func(err error, ret any) {
		if err != nil {
			l.Log.Infof("chatremote.leave ret, err: %v", err)
			return
		}
		l.Log.Infof("chatremote.leave ret")
	})
}
