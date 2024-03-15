package client

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	client "github.com/dfklegend/cell2/pomeloclient"
	"github.com/dfklegend/cell2/pomelonet/common/conn/message"
	"github.com/dfklegend/cell2/utils/common"
	"github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/runservice"

	"client/protos"
)

const (
	STATE_INIT = iota
	STATE_CONNECT_GATE
)

type ChatClient struct {
	//
	Client     *client.CellClient
	RunService *runservice.StandardRunService
	state      int
	msgBind    bool
	Username   string
}

func NewChatClient(name string) *ChatClient {
	c := client.NewCellClient(name, nil)
	r := c.GetRunService()

	return &ChatClient{
		Client:     c,
		RunService: r,
		state:      STATE_INIT,
		msgBind:    false,
		Username:   name,
	}
}

func (s *ChatClient) setState(state int) {
	s.state = state
}

func (s *ChatClient) getState() int {
	return s.state
}

func (s *ChatClient) Start(address string) {
	c := s.Client
	c.Start(address)

	c.SetCBConnected(s.onConnected)
	c.SetCBBreak(func() {
		s.setState(STATE_INIT)
	})
	s.bindMsgs()
}

func (s *ChatClient) Stop() {
	s.Client.Stop()
}

func (s *ChatClient) onConnected() {
	logger.Log.Debugf("onConnected")
	if s.getState() != STATE_INIT {
		return
	}

	// 开始
	s.queryGate()
	s.setState(STATE_CONNECT_GATE)
}

func (s *ChatClient) bindMsgs() {
	if s.msgBind {
		return
	}
	s.msgBind = true
	// 注册
	//ec := s.RunService.GetEventCenter()
}

func (s *ChatClient) queryGate() {
	logger.Log.Debugf("queryGate")

	c := s.Client
	c.GetClient().SendRequest("handler.handler.querygate", []byte("{}"), func(error bool, msg *message.Message) {
		fmt.Println("ack from cb")
		fmt.Println(string(msg.Data))

		if error {
			fmt.Println("query handler error ")
			return
		}

		// TODO:解析，使用正确端口
		port := s.getGatePort(msg.Data)
		c.Start(fmt.Sprintf("127.0.0.1:%v", port))
		c.WaitReady()

		m := &protos.LoginReq{
			Username: s.Username,
		}
		c.GetClient().SendRequest("handler.handler.login", []byte(common.SafeJsonMarshal(m)), s.onLoginRet)
	})
}

func (s *ChatClient) getGatePort(data []byte) string {
	obj := protos.QueryGateAck{}
	json.Unmarshal(data, &obj)

	port := obj.Port
	subs := strings.Split(port, ",")
	if len(subs) >= 2 {
		return subs[1]
	}
	return ""
}

func (s *ChatClient) onLoginRet(error bool, msg *message.Message) {
	if error {
		logger.L.Errorf("login failed!")
		return
	}
	// 可以发送消息
	obj := protos.LoginAck{}
	json.Unmarshal(msg.Data, &obj)
	logger.Log.Infof("onLoginRet: %+v", obj)

	if obj.Code == 0 {
		s.Client.GetRunService().GetTimerMgr().AddTimer(5*time.Second, func(args ...interface{}) {
			s.startGame()
		})
	}
}

func (s *ChatClient) startGame() {
	m := &protos.EmptyArgReq{}
	s.Client.GetClient().SendRequest("logic_old.logic_old.startfight", []byte(common.SafeJsonMarshal(m)), func(err bool, msg *message.Message) {
		if err {
			logger.L.Errorf("start game error")
		}
	})
}
