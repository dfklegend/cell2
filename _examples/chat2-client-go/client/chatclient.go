package client

import (
	"fmt"
	"strings"

	"github.com/dfklegend/cell2/apimapper/apientry"
	client "github.com/dfklegend/cell2/pomeloclient"
	"github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/runservice"
	jsonserializer "github.com/dfklegend/cell2/utils/serialize/json"

	"chat2-client-go/protos"
)

const (
	STATE_INIT = iota
	STATE_CONNECT_GATE
)

type ChatClient struct {
	Client     *client.CellClient
	RunService *runservice.StandardRunService
	state      int
	msgBind    bool
}

func NewChatClient(name string) *ChatClient {
	c := client.NewCellClient(name, jsonserializer.GetDefaultSerializer())
	r := c.GetRunService()

	cc := &ChatClient{
		Client:     c,
		RunService: r,
		state:      STATE_INIT,
		msgBind:    false,
	}
	c.SetOwner(cc)
	c.RegisterDefaultEntry(&Entry{}, apientry.WithNameFunc(apientry.ToLowerCamelCase))
	c.BuildEntries()
	return cc
}

func (self *ChatClient) setState(s int) {
	self.state = s
}

func (self *ChatClient) getState() int {
	return self.state
}

func (self *ChatClient) Start(address string) {
	c := self.Client
	c.Start(address)

	c.SetCBConnected(self.onConnected)
	c.SetCBBreak(func() {
		self.setState(STATE_INIT)
	})
}

func (self *ChatClient) Stop() {
	self.Client.Stop()
}

func (self *ChatClient) onConnected() {
	logger.Log.Debugf("onConnected")
	if self.getState() != STATE_INIT {
		return
	}

	// 开始
	self.queryGate()
	self.setState(STATE_CONNECT_GATE)
}

func (self *ChatClient) queryGate() {
	logger.Log.Debugf("queryGate")

	c := self.Client
	c.SendRequest("gate.gate.querygate", nil, func(err error, ret interface{}) {
		ack := ret.(*protos.QueryGateAck)
		self.onQueryGateAck(ack)
	}, TypePtrQueryGateAck)
}

func (self *ChatClient) onQueryGateAck(ack *protos.QueryGateAck) {
	port := self.getGatePort(ack)
	c := self.Client
	c.Start(fmt.Sprintf("127.0.0.1:%v", port))
	c.WaitReady()

	c.SendRequest("gate.gate.login", &protos.LoginReq{
		Name: "fromGoClient",
	}, self.onLoginRet, TypePtrNormalAck)
}

func (self *ChatClient) getGatePort(ack *protos.QueryGateAck) string {
	port := ack.Port
	subs := strings.Split(port, ",")
	if len(subs) >= 2 {
		return subs[1]
	}
	return ""
}

func (self *ChatClient) onLoginRet(err error, ret interface{}) {
	// 可以发送消息
	obj := ret.(*protos.NormalAck)
	logger.Log.Debugf("onLoginRet: %+v", obj)

	c := self.Client

	c.SendRequest("chat.chat.sendchat", &protos.ChatMsg{
		Name:    "fromGoClient",
		Content: "hello",
	}, func(err error, ret interface{}) {
		if err != nil {
			logger.Log.Infof("send chat ret, error:%v", err)
		} else {
			logger.Log.Infof("send chat succ")
		}
	}, TypePtrNormalAck)

	//c.GetClient().SendNotify("chat.chat.sendchat1", data)
	//c.GetClient().SendRequest("chat.chat.sendchat1", data, nil)
	//c.GetClient().SendNotify("chat.chat.sendchat", data)
}
