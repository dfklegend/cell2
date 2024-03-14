// Copyright (c) TFG Co. All Rights Reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dfklegend/cell2/pomelonet/common/conn/codec"
	"github.com/dfklegend/cell2/pomelonet/common/conn/message"
	"github.com/dfklegend/cell2/pomelonet/common/conn/packet"
	"github.com/dfklegend/cell2/utils/compression"
	"github.com/dfklegend/cell2/utils/logger"
)

const (
	DefaultReqTimeout   = 15 * time.Second
	DefaultWriteTimeout = 15 * time.Second
	MaxQueueSize        = 999
)

var (
	handshakeBuffer = `
{
	"sys": {
		"platform": "mac",
		"libVersion": "0.3.5-release",
		"clientBuildNumber":"20",
		"clientVersion":"2.1"
	},
	"user": {
		"age": 30
	}
}
`
)

// HandshakeSys struct
type HandshakeSys struct {
	Dict       map[string]uint16 `json:"dict"`
	Heartbeat  int               `json:"heartbeat"`
	Serializer string            `json:"serializer"`
}

// HandshakeData struct
type HandshakeData struct {
	Code int          `json:"code"`
	Sys  HandshakeSys `json:"sys"`
}

type ClientMsg struct {
	Msg *message.Message
	Cb  CBTask
}

// 加入一个callback
type pendingRequest struct {
	msg    *message.Message
	sentAt time.Time
	cb     CBTask
}

type CBTask func(error bool, msg *message.Message)

// Client
// 重连在上层处理
// 携程
//     readServerMessages
//     		读取消息
//     		组织成packet
//     handlePackets
//     		处理packet
//     sendHeartbeats
//     		定时发送heartbeat
//     pendingRequestsReaper
//     		移除超时未返回的请求
type Client struct {
	name      string
	conn      net.Conn
	Connected bool
	// 握手完毕
	Ready           bool
	packetEncoder   codec.PacketEncoder
	packetDecoder   codec.PacketDecoder
	packetChan      chan *packet.Packet
	IncomingMsgChan chan *ClientMsg
	// 和pendingReqMutex 会造成死锁
	//pendingChan     chan bool
	pendingRequests map[uint]*pendingRequest
	pendingReqMutex sync.Mutex
	requestTimeout  time.Duration

	sendQueue      chan []byte
	closeChan      chan struct{}
	nextID         uint32
	messageEncoder message.Encoder

	released  int32
	detailLog bool
}

// MsgChannel return the incoming message channel
func (c *Client) MsgChannel() chan *ClientMsg {
	return c.IncomingMsgChan
}

// ConnectedStatus return the connection status
func (c *Client) ConnectedStatus() bool {
	return c.Connected
}

// New returns a new client
func New(requestTimeout ...time.Duration) *Client {

	reqTimeout := DefaultReqTimeout
	if len(requestTimeout) > 0 {
		reqTimeout = requestTimeout[0]
	}

	return &Client{
		name:          "",
		Connected:     false,
		Ready:         false,
		packetEncoder: codec.NewPomeloPacketEncoder(),
		packetDecoder: codec.NewPomeloPacketDecoder(),

		packetChan:      make(chan *packet.Packet, 10),
		IncomingMsgChan: make(chan *ClientMsg, MaxQueueSize),

		pendingRequests: make(map[uint]*pendingRequest),
		sendQueue:       make(chan []byte, MaxQueueSize),
		requestTimeout:  reqTimeout,
		// 30 here is the limit of inflight messages
		// TODO this should probably be configurable
		//pendingChan:    make(chan bool, 30),
		messageEncoder: message.NewMessagesEncoder(false),
		released:       0,
		detailLog:      false,
	}
}

func (c *Client) SetName(name string) {
	c.name = name
}

func (c *Client) SetDetailLog(v bool) {
	c.detailLog = v
}

func (c *Client) onError() {
}

// Disconnect disconnects the client
// can not use again
func (c *Client) Disconnect() {
	if !atomic.CompareAndSwapInt32(&c.released, 0, 1) {
		return
	}

	if c.Connected {
		c.Connected = false
		c.Ready = false

		close(c.closeChan)
		c.conn.Close()
	}
	close(c.IncomingMsgChan)
	close(c.packetChan)
	close(c.sendQueue)
	c.packetChan = nil
}

// ConnectToTLS connects to the server at addr using TLS, for now the only supported protocol is tcp
// this methods blocks as it also handles the messages from the server
func (c *Client) ConnectToTLS(addr string, skipVerify bool) error {
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: skipVerify,
	})
	if err != nil {
		return err
	}
	return c.start(conn)
}

// ConnectTo connects to the server at addr, for now the only supported protocol is tcp
// this methods blocks as it also handles the messages from the server
func (c *Client) ConnectTo(addr string) error {
	fmt.Printf("begin dial:%v\n", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("dial error:%v\n", err)
		return err
	}
	return c.start(conn)
}

func (c *Client) start(conn net.Conn) error {
	c.conn = conn

	fmt.Printf("dial over\n")

	go c.doSend()

	// 等一会
	//fmt.Println("wait 5 seconds")
	//time.Sleep(5 * time.Second)
	if err := c.beginHandshake(); err != nil {
		return err
	}

	c.closeChan = make(chan struct{})

	return nil
}

func (c *Client) sendHandshakeRequest() error {
	p, err := c.packetEncoder.Encode(packet.Handshake, []byte(handshakeBuffer))
	if err != nil {
		return err
	}
	_, err = c.conn.Write(p)
	return err
}

func (c *Client) beginRecvMessages() error {
	c.Connected = true

	go c.readServerMessages()
	go c.handlePackets()
	return nil
}

// 等到一个握手返回，整个连接就建立成功
func (c *Client) handleHandleShake(handshakePacket *packet.Packet) error {
	fmt.Println("got handleHandleShake")
	if handshakePacket.Type != packet.Handshake {
		return fmt.Errorf("got first packet from server that is not a handshake, aborting")
	}

	var err error
	handshake := &HandshakeData{}
	if compression.IsCompressed(handshakePacket.Data) {
		handshakePacket.Data, err = compression.InflateData(handshakePacket.Data)
		if err != nil {
			return err
		}
	}

	err = json.Unmarshal(handshakePacket.Data, handshake)
	if err != nil {
		return err
	}

	logger.Log.Debug("got handshake from sv, data: %v", handshake)

	if handshake.Sys.Dict != nil {
		message.SetDictionary(handshake.Sys.Dict)
	}
	p, err := c.packetEncoder.Encode(packet.HandshakeAck, []byte{})
	if err != nil {
		return err
	}
	_, err = c.conn.Write(p)
	if err != nil {
		return err
	}

	logger.Log.Debug("connection ready")
	// 握手成功
	c.Ready = true

	// 握手成功
	go c.sendHeartbeats(handshake.Sys.Heartbeat)
	go c.pendingRequestsReaper()
	return nil
}

// pendingRequestsReaper delete timedout requests
func (c *Client) pendingRequestsReaper() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	defer func() {
		if err := recover(); err != nil {
			// maybe closed
		}
	}()

	for {
		select {
		case <-ticker.C:
			if c.detailLog {
				logger.Log.Infof("%v checktimeout pendingnum: %v incomemsgnum: %v",
					c.name, len(c.pendingRequests), len(c.IncomingMsgChan))
			}

			toDelete := make([]*pendingRequest, 0)

			//logger.Log.Infof("pendingRequestsReaper pending lock")
			c.pendingReqMutex.Lock()
			//logger.Log.Infof("pendingRequestsReaper pending locked")
			num := 0
			for _, v := range c.pendingRequests {
				if time.Now().Sub(v.sentAt) > c.requestTimeout {
					toDelete = append(toDelete, v)
					num++
				}
				if num > 100 {
					break
				}
			}
			c.pendingReqMutex.Unlock()

			for _, pendingReq := range toDelete {
				//err := errors.New("request timeout")

				// send a err msg to incoming msg chan
				m := &message.Message{
					Type:  message.Response,
					ID:    pendingReq.msg.ID,
					Route: pendingReq.msg.Route,
					Data:  nil,
					Err:   true,
				}
				nm := &ClientMsg{
					Msg: m,
					Cb:  pendingReq.cb,
				}

				logger.Log.Warnf("request timeout: %v", pendingReq.msg.ID)

				c.pendingReqMutex.Lock()
				delete(c.pendingRequests, pendingReq.msg.ID)
				c.pendingReqMutex.Unlock()

				//logger.Log.Infof("pendingRequestsReaper -> in")
				//<-c.pendingChan
				//logger.Log.Infof("pendingRequestsReaper -> out")

				//logger.Log.Infof("request timeout, :%v", m.Route)
				c.pushIncomingMsg(nm)
			}

			if c.detailLog && len(toDelete) > 0 {
				logger.Log.Warnf("%v push timeout num: %v pendingnum: %v incomemsgnum: %v",
					c.name, len(toDelete), len(c.pendingRequests), len(c.IncomingMsgChan))
			}

			if c.detailLog {
				logger.Log.Infof("%v checktimeout over pendingnum: %v incomemsgnum: %v",
					c.name, len(c.pendingRequests), len(c.IncomingMsgChan))
			}

			//logger.Log.Infof("pendingRequestsReaper pending unlocked")
		case <-c.closeChan:
			return
		}
	}
}

func (c *Client) pushIncomingMsg(msg *ClientMsg) {
	defer func() {
		if err := recover(); err != nil {
			// may be closed
		}
	}()

	if c.detailLog {
		if msg.Msg.ID > 0 {
			logger.Log.Infof("%v got response ID: %v", c.name, msg.Msg.ID)
		}

	}
	c.IncomingMsgChan <- msg
}

// 处理包
func (c *Client) handlePackets() {
	for {
		select {
		case p := <-c.packetChan:
			switch p.Type {
			case packet.Handshake:
				c.handleHandleShake(p)
			case packet.Data:
				//handle data
				//logger.Log.Debug("got data: %s", string(p.Data))
				var cbTask CBTask
				m, err := message.Decode(p.Data)
				if err != nil {
					logger.Log.Errorf("error decoding msg from sv: %s", string(m.Data))
				}
				if m.Type == message.Response {
					//logger.Log.Infof("handlePackets pending lock")
					c.pendingReqMutex.Lock()
					//logger.Log.Infof("handlePackets pending locked")

					if req, ok := c.pendingRequests[m.ID]; ok {
						cbTask = req.cb
						delete(c.pendingRequests, m.ID)
						//logger.Log.Infof("handpanckets -> in")
						//<-c.pendingChan
						//logger.Log.Infof("handpanckets -> out")
					} else {
						// 扔掉
						c.pendingReqMutex.Unlock()
						logger.Log.Warnf("drop response: %v", m.ID)
						continue // do not process msg for already timedout request
					}
					c.pendingReqMutex.Unlock()
					//logger.Log.Infof("handlePackets pending unlocked")
				}
				//fmt.Println("cbTask")
				//fmt.Println(cbTask)
				//if cbTask != nil {
				//	cbTask(false, m)
				//}
				nm := &ClientMsg{
					Msg: m,
					Cb:  cbTask,
				}
				c.pushIncomingMsg(nm)
			case packet.Kick:
				logger.Log.Warn("got kick packet from the server! disconnecting...")
				c.Disconnect()
			}
		case <-c.closeChan:
			return
		}
	}
}

func (c *Client) readPackets(buf *bytes.Buffer) ([]*packet.Packet, error) {
	// listen for sv messages
	data := make([]byte, 1024)
	n := len(data)
	var err error

	for n == len(data) {
		n, err = c.conn.Read(data)
		//fmt.Printf("conn.read ret:%d %v", n, err)
		if err != nil {
			return nil, err
		}
		buf.Write(data[:n])
	}
	packets, err := c.packetDecoder.Decode(buf.Bytes())
	if err != nil {
		logger.Log.Errorf("error decoding packet from server: %s", err.Error())
	}
	totalProcessed := 0
	for _, p := range packets {
		totalProcessed += codec.HeadLength + p.Length
	}
	buf.Next(totalProcessed)

	return packets, nil
}

// 读取包
func (c *Client) readServerMessages() {
	defer c.Disconnect()

	buf := bytes.NewBuffer(nil)
	for c.Connected {
		packets, err := c.readPackets(buf)
		if err != nil && c.Connected {
			logger.Log.Errorf("%v read error: %v", c.name, err)
			break
		}

		for _, p := range packets {
			c.pushPacket(p)
		}
	}
}

func (c *Client) pushPacket(p *packet.Packet) {
	defer func() {
		if err := recover(); err != nil {
			// maybe closed
		}
	}()
	c.packetChan <- p
}

func (c *Client) sendHeartbeats(interval int) {
	t := time.NewTicker(time.Duration(interval) * time.Second)
	defer func() {
		t.Stop()
		c.Disconnect()
	}()
	for {
		select {
		case <-t.C:
			p, _ := c.packetEncoder.Encode(packet.Heartbeat, []byte{})

			//c.conn.SetWriteDeadline(time.Now().Add(DefaultWriteTimeout))
			////logger.Log.Infof("sendHeartbeats size: %v", len(p))
			//_, err := c.conn.Write(p)
			//if err != nil {
			//	logger.Log.Errorf("error sending heartbeat to server: %s", err.Error())
			//	return
			//}
			c.safePushSend(p)
		case <-c.closeChan:
			return
		}
	}
}

func (c *Client) beginHandshake() error {
	if err := c.sendHandshakeRequest(); err != nil {
		return err
	}

	if err := c.beginRecvMessages(); err != nil {
		return err
	}
	return nil
}

// SendRequest sends a request to the server
func (c *Client) SendRequest(route string, data []byte, cb CBTask) (uint, error) {
	return c.sendMsg(message.Request, route, data, cb)
}

// SendNotify sends a notify to the server
func (c *Client) SendNotify(route string, data []byte) error {
	_, err := c.sendMsg(message.Notify, route, data, nil)
	return err
}

func (c *Client) buildPacket(msg message.Message) ([]byte, error) {
	encMsg, err := c.messageEncoder.Encode(&msg)
	if err != nil {
		return nil, err
	}
	p, err := c.packetEncoder.Encode(packet.Data, encMsg)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// sendMsg sends the request to the server
func (c *Client) sendMsg(msgType message.Type, route string, data []byte, cb CBTask) (uint, error) {
	//
	if !c.Ready {
		logger.Log.Errorf("sendMsg error, client not Ready: %v", route)
		return 0, nil
	}

	m := message.Message{
		Type:  msgType,
		ID:    uint(atomic.AddUint32(&c.nextID, 1)),
		Route: route,
		Data:  data,
		Err:   false,
	}
	p, err := c.buildPacket(m)
	if msgType == message.Request {

		//logger.Log.Infof("sendMsg <- in")
		//c.pendingChan <- true
		//logger.Log.Infof("sendMsg <- out")

		//logger.Log.Infof("sendMsg pending lock")
		c.pendingReqMutex.Lock()
		//logger.Log.Infof("sendMsg pending locked")
		if _, ok := c.pendingRequests[m.ID]; !ok {
			c.pendingRequests[m.ID] = &pendingRequest{
				msg:    &m,
				sentAt: time.Now(),
				cb:     cb,
			}
		}
		c.pendingReqMutex.Unlock()
		//logger.Log.Infof("sendMsg pending unlocked")
	}

	if err != nil {
		return m.ID, err
	}

	//fmt.Printf("conn:%v\n", c.conn)
	//logger.Log.Infof("sendMsg size: %v", len(p))

	//c.conn.SetWriteDeadline(time.Now().Add(DefaultWriteTimeout))
	//_, err = c.conn.Write(p)
	//if err != nil {
	//	logger.Log.Errorf("Write error: %v", err)
	//	c.onError()
	//}
	//return m.ID, err

	if c.detailLog {
		if msgType == message.Request {
			logger.Log.Infof("%v push request ID: %v", c.name, m.ID)
		}

	}
	c.safePushSend(p)

	if c.detailLog {
		logger.Log.Infof("%v sendqueue size: %v", c.name, len(c.sendQueue))
	}
	return m.ID, nil
}

func (c *Client) safePushSend(p []byte) {
	defer func() {
		if err := recover(); err != nil {
			// maybe closed
		}
	}()
	c.sendQueue <- p
}

func (c *Client) doSend() {
	// 可以一次多获取一些buf来发送
	for {
		select {
		case buf := <-c.sendQueue:
			{
				// 设置发送超时，在cpu压力大的情况下，会导致很多发送被取消了
				//c.conn.SetWriteDeadline(time.Now().Add(DefaultWriteTimeout))
				sendSize, err := c.conn.Write(buf)
				if err != nil {
					logger.Log.Errorf("%v Write error: %v", c.name, err)
					c.onError()
				}
				if sendSize != len(buf) {
					logger.Log.Errorf("Write size error: %v != %v", sendSize, len(buf))
				}
			}
		case <-c.closeChan:
			return
		}
	}
}
