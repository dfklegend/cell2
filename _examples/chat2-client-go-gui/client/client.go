package client

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dfklegend/cell2/utils/jsonutils"

	pomelo "github.com/revzim/go-pomelo-client"

	"chat2-client-go-gui/protos"
)

//NetTick 网络帧率
const NetTick = 10

//Heartbeat 心跳
const Heartbeat = 30

func newClient(name string) *Client {
	return &Client{
		PomeloClient: pomelo.NewConnector(),
		Name:         name,
	}
}

type Client struct {
	PomeloClient *pomelo.Connector
	Name         string

	Queue *[]func()
}

func start(host string) {
	if cfg.Client != nil {
		cfg.Client.close()
	}

	c := newClient(cfg.UserName)
	cfg.Client = c
	c.registerHandler()
	cfg.Host = host

	var err error
	err = c.PomeloClient.InitReqHandshake("0.6.0", "golang-pomelo", nil, map[string]interface{}{
		"age": 30,
	})
	checkError(err)

	err = c.PomeloClient.InitHandshakeACK(Heartbeat)
	checkError(err)

	c.PomeloClient.Connected(func() {
		log.Printf("client connected : %s", host)
		if cfg.newGateway {
			c.login()
			return
		}
		cfg.newGateway = true
		c.queryGate()
	})

	go func() {
		err = c.PomeloClient.Run(host, false, int64(NetTick))
		log.Println(err)
	}()
}

func (c *Client) registerHandler() {
	c.PomeloClient.On("onNewUser", c.onNewUser)
	c.PomeloClient.On("onMessage", c.onMessage)
	c.PomeloClient.On("onMembers", c.onMembers)
	c.PomeloClient.On("onUserLeave", c.onUserLeave)
}

func (c *Client) queryGate() {
	var msg protos.EmptyArgReq
	c.Request("gate.gate.querygate", msg, c.onQueryGate)
}

func (c *Client) onQueryGate(data []byte) {
	msg := protos.QueryGateAck{}

	jsonutils.Unmarshal(data, &msg)

	//获取tcp端口
	//{0 127.0.0.1 30012,30022}
	port := strings.Split(msg.Port, ",")[1]

	start(fmt.Sprintf("%s:%s", msg.IP, port))
}

func (c *Client) login() {
	msg := protos.LoginReq{
		Name: "tong",
	}
	c.Request("gate.gate.login", msg, c.onLogin)
}

func (c *Client) onLogin(data []byte) {
	msg := protos.NormalAck{}
	jsonutils.Unmarshal(data, &msg)
	s := jsonutils.Marshal(msg)
	log.Printf("[recv] onLogin : %v", s)

	cfg.onSystem(fmt.Sprintf("[login] : %v", s))
	cfg.Login = true
}

func (c *Client) onNewUser(data []byte) {
	msg := protos.OnNewUser{}
	jsonutils.Unmarshal(data, &msg)
	s := jsonutils.Marshal(msg)
	log.Printf("[recv] onNewUser : %v", s)

	cfg.onSystem(fmt.Sprintf("onNewUser %s", msg.Name))
}

func (c *Client) onMessage(data []byte) {
	msg := protos.ChatMsg{}
	jsonutils.Unmarshal(data, &msg)
	s := jsonutils.Marshal(msg)
	log.Printf("[recv] onMessage : %v", s)

	cfg.onMessage(msg.Content, msg.Name)
}

func (c *Client) onMembers(data []byte) {
	msg := protos.OnMembers{}
	jsonutils.Unmarshal(data, &msg)
	s := jsonutils.Marshal(msg)
	log.Printf("[recv] OnMembers : %v", s)
	cfg.onSystem(fmt.Sprintf("OnMembers %s", msg.Members))
}

func (c *Client) onUserLeave(data []byte) {
	msg := protos.OnUserLeave{}
	jsonutils.Unmarshal(data, &msg)
	s := jsonutils.Marshal(msg)
	log.Printf("[recv] onUserLeave : %v", s)
	cfg.onSystem(fmt.Sprintf("onUserLeave %s", msg.Name))
}

func (c *Client) Notify(route string, msg any) {
	c.PomeloClient.Notify(route, []byte(jsonutils.Marshal(&msg)))
}

func (c *Client) Request(route string, msg any, cb func(data []byte)) {
	c.PomeloClient.Request(route, []byte(jsonutils.Marshal(&msg)), cb)

}

func (c *Client) sendChatMsg(content string) {
	m := protos.ChatMsg{
		Name:    c.Name,
		Content: content,
	}
	c.Notify("chat.chat.sendchat", m)
}

func (c *Client) close() {
	if c.PomeloClient.IsClosed() {
		return
	}
	c.PomeloClient.Close()
}

func checkError(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

var cfg Cfg

type Cfg struct {
	Host     string
	UserName string
	ChatMsg  string
	Room     string

	ChatHistory []ChatMsg
	Client      *Client
	newGateway  bool
	Login       bool
}

func (c *Cfg) onSystem(msg string) {
	c.ChatHistory = append(c.ChatHistory, ChatMsg{
		Time: time.Now(),
		Name: "system",
		Msg:  msg,
	})
}

func (c *Cfg) onMessage(msg string, name string) {
	c.ChatHistory = append(c.ChatHistory, ChatMsg{
		Time: time.Now(),
		Name: name,
		Msg:  msg,
	})

}

type ChatMsg struct {
	Time time.Time
	Name string
	Msg  string
}
