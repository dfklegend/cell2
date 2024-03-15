package telnetcmd

import (
	"fmt"
	"reflect"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/reiver/go-telnet"

	"github.com/dfklegend/cell2/actorex"
	"github.com/dfklegend/cell2/actorex/service"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/utils/sche"
	"master/messages"
)

const (
	DefaultTelnetPort = 5555
)

// 获取控制台输入
// 转发给servicecmd
// 获取输出并显示

type Service struct {
	*service.Service

	handler *Handler

	cmdService *actor.PID
}

func NewService() *Service {
	s := &Service{
		Service: service.NewService(),
		handler: NewTelnetHandler(),
	}

	s.Service.InitReqReceiver(s)
	return s
}

func (s *Service) Start() {
	s.cmdService = actor.NewPID(actorex.LocalAddress, "cmdservice")

	go s.startTelnetServer()
	s.addInputProcess()
}

func (s *Service) startTelnetServer() {
	port := app.Node.GetMasterInfo().TelnetPort
	if port == 0 {
		port = DefaultTelnetPort
	}

	server := &telnet.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: s.handler,
	}

	err := server.ListenAndServe()
	if nil != err {
		panic(err)
	}
}

func (s *Service) addInputProcess() {
	rs := s.GetRunService()
	rs.GetSelector().AddSelector("inputprocess", sche.NewFuncSelector(reflect.ValueOf(s.handler.Input),
		func(v reflect.Value, recvOk bool) {
			if !recvOk {
				return
			}

			input := v.Interface().(*sessionInput)

			//input.Session.pushOutput(input.Str)
			s.processInput(input.Session, input.Str)
		}))
}

func (s *Service) processInput(session *session, cmd string) {
	// 请求service cmd
	s.RequestEx(s.cmdService, "service.cmd", &messages.Cmd{
		Cmd: cmd,
	}, func(err error, raw interface{}) {
		if err != nil {
			return
		}
		ack := raw.(*messages.CmdAck)

		//fmt.Println(ack.Result)
		session.pushOutput(ack.Result)
	})
}
