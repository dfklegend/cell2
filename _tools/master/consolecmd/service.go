package consolecmd

import (
	"bufio"
	"fmt"
	"os"
	"reflect"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/actorex"
	"github.com/dfklegend/cell2/actorex/service"
	"github.com/dfklegend/cell2/utils/sche"
	"master/messages"
)

// 获取控制台输入
// 转发给servicecmd
// 获取输出并显示

type Service struct {
	*service.Service

	input      chan string
	cmdService *actor.PID
}

func NewService() *Service {
	s := &Service{
		Service: service.NewService(),
		input:   make(chan string, 9),
	}

	s.Service.InitReqReceiver(s)
	return s
}

func (s *Service) Start() {
	s.cmdService = actor.NewPID(actorex.LocalAddress, "cmdservice")
	s.addInputProcess()
	go s.loopInput()
}

func (s *Service) loopInput() {
	in := bufio.NewReader(os.Stdin)
	for true {
		str, _, _ := in.ReadLine()
		if len(str) == 0 {
			fmt.Print("please input:")
		} else {
			s.input <- string(str)
		}
	}
}

func (s *Service) addInputProcess() {
	rs := s.GetRunService()
	rs.GetSelector().AddSelector("inputprocess", sche.NewFuncSelector(reflect.ValueOf(s.input),
		func(v reflect.Value, recvOk bool) {
			if !recvOk {
				return
			}

			str := v.Interface().(string)
			s.processInput(str)
		}))
}

func (s *Service) processInput(cmd string) {
	// 请求service cmd
	//fmt.Printf("cmd : %v\n", cmd)
	// send it to cmdservice
	s.RequestEx(s.cmdService, "service.cmd", &messages.Cmd{
		Cmd: cmd,
	}, func(err error, raw interface{}) {
		if err != nil {
			return
		}
		ack := raw.(*messages.CmdAck)

		fmt.Println(ack.Result)
	})
}
