// Copyright (c) nano Author and TFG Co. All Rights Reserved.
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

package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dfklegend/cell2/utils/common"
	"github.com/dfklegend/cell2/utils/compression"
	"github.com/dfklegend/cell2/utils/logger"

	"github.com/dfklegend/cell2/pomelonet/common/conn/codec"
	"github.com/dfklegend/cell2/pomelonet/common/conn/message"
	"github.com/dfklegend/cell2/pomelonet/common/conn/packet"
	"github.com/dfklegend/cell2/pomelonet/constants"
	"github.com/dfklegend/cell2/pomelonet/interfaces"
	"github.com/dfklegend/cell2/pomelonet/server/acceptor"
	//"github.com/dfklegend/cell2/server/interfaces"
)

// session state
const (
	_ int32 = iota
	// StatusStart status
	StatusStart
	// StatusHandshake status
	StatusHandshake
	// StatusWorking status
	StatusWorking
	// StatusClosed status
	StatusClosed
)

var (
	// hbd contains the heartbeat packet data
	hbd []byte
	// hrd contains the handshake response data
	hrd  []byte
	once sync.Once

	DefaultHeartbeatTimeSeconds int64 = 10
	DefaultHeartbeatTime              = 10 * time.Second

	nilBytes = make([]byte, 0)
)

type pendingWrite struct {
	data []byte
	err  error
}

type pendingMessage struct {
	typ     message.Type // message type
	route   string       // message route (push)
	mid     uint         // response message id (response)
	payload interface{}  // payload
	err     bool         // if its an error message
}

//
type ISession interface {
	Handle()
}

// ClientSession
// 负责接收来自客户端的消息
type ClientSession struct {
	// set after add to frontSessions
	netId uint32
	// 连接对象
	conn acceptor.PlayerConn
	// 发送队列
	chSend chan *pendingWrite

	cfg *SessionConfig
	// 注:目前是nil
	// 因为下发数据都是已经序列化过了[]byte
	//serializer      serialize.Serializer

	chanClose chan bool
	mutex     sync.Mutex

	state         int32 // current agent state
	lastHeartBeat int64
}

func NewClientSession(
	c acceptor.PlayerConn,
	cfg *SessionConfig) *ClientSession {

	once.Do(func() {
		serializerName := "json"
		heartbeatTime := DefaultHeartbeatTime
		hbdEncode(heartbeatTime, cfg.Encoder,
			cfg.MessageEncoder.IsCompressionEnabled(), serializerName)
	})

	v := &ClientSession{
		conn:      c,
		cfg:       cfg,
		chSend:    make(chan *pendingWrite, 9999),
		chanClose: make(chan bool),
		state:     StatusStart,
	}

	// TODO
	//GetFrontSessions().AddSession(v)
	v.GetImpl().OnSessionCreate(v)
	return v
}

func (s *ClientSession) GetImpl() interfaces.IClientSessionImpl {
	return s.cfg.Impl
}

func (s *ClientSession) Reserve() {}

// 由管理器传入
func (s *ClientSession) SetId(id uint32) {
	s.netId = id
}

func (s *ClientSession) GetId() uint32 {
	return s.netId
}

func (s *ClientSession) SetStatus(state int32) {
	atomic.StoreInt32(&s.state, state)
}

// GetStatus gets the status
func (s *ClientSession) GetStatus() int32 {
	return atomic.LoadInt32(&s.state)
}

func (s *ClientSession) IsClosed() bool {
	return s.GetStatus() == StatusClosed
}

func (s *ClientSession) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	select {
	// close already
	case <-s.chanClose:
		return
	default:
		s.SetStatus(StatusClosed)
		close(s.chanClose)
		close(s.chSend)
	}

	logger.L.Infof("clientsession close %v", s.netId)

	// close channel
	s.conn.Close()
	s.GetImpl().OnSessionClose(s)
}

// 读写消息
func (s *ClientSession) Handle() {
	go s.heartbeat()
	go s.write()
	go s.read()
}

// ---- write begin ----
func (s *ClientSession) write() {
	// clean func
	defer func() {
		s.Close()
	}()

	for {
		select {
		case <-s.chanClose:
			return
		case pWrite, ok := <-s.chSend:
			if !ok {
				return
			}
			// pWrite中的错误信息不会发到客户端
			// close agent if low-level Conn broken
			if _, err := s.conn.Write(pWrite.data); err != nil {
				logger.Log.Errorf("Failed to write in conn: %s", err.Error())
				return
			}
		}
	}
}

func (s *ClientSession) getMessageFromPendingMessage(pm pendingMessage) (*message.Message, error) {
	var payload []byte
	var ok bool
	if pm.payload == nil {
		payload = nilBytes
	} else {
		payload, ok = pm.payload.([]byte)
		if !ok {
			// TODO: 错误
			return nil, errors.New("bad payload")
		}
	}

	// construct message and encode
	m := &message.Message{
		Type:  pm.typ,
		Data:  payload,
		Route: pm.route,
		ID:    pm.mid,
		Err:   pm.err,
	}

	return m, nil
}

func (s *ClientSession) send(pendingMsg pendingMessage) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("send error")
		}
	}()

	m, err := s.getMessageFromPendingMessage(pendingMsg)
	if err != nil {
		return err
	}

	// packet encode
	p, err := s.packetEncodeMessage(m)
	if err != nil {
		return err
	}

	pWrite := &pendingWrite{
		data: p,
	}

	if pendingMsg.err {
		pWrite.err = errors.New("getMessageFromPendingMessage error")
	}

	s.pushToSend(pWrite)
	return
}

func (s *ClientSession) packetEncodeMessage(m *message.Message) ([]byte, error) {
	em, err := s.cfg.MessageEncoder.Encode(m)
	if err != nil {
		return nil, err
	}

	// packet encode
	p, err := s.cfg.Encoder.Encode(packet.Data, em)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// Push implementation for NetworkEntity interface
func (s *ClientSession) Push(route string, v interface{}) error {
	if s.GetStatus() == StatusClosed {
		return errors.New("closed")
	}

	// switch d := v.(type) {
	// case []byte:
	//     logger.Log.Debugf("Type=Push, ID=%d, UID=%s, Route=%s, Data=%dbytes",
	//         s.Session.ID(), a.Session.UID(), route, len(d))
	// default:
	//     logger.Log.Debugf("Type=Push, ID=%d, UID=%s, Route=%s, Data=%+v",
	//         a.Session.ID(), a.Session.UID(), route, v)
	// }
	return s.send(pendingMessage{typ: message.Push, route: route, payload: v})
}

// SendHandshakeResponse sends a handshake response
func (s *ClientSession) SendHandshakeResponse() error {
	_, err := s.conn.Write(hrd)
	return err
}

// ResponseMID implementation for NetworkEntity interface
// Respond message to session
func (s *ClientSession) ResponseMID(mid uint, v interface{}, e error) error {
	err := false
	if e != nil {
		err = true
	}
	if s.GetStatus() == StatusClosed {
		return errors.New("closed")
	}

	if mid <= 0 {
		return constants.ErrSessionOnNotify
	}

	// switch d := v.(type) {
	// case []byte:
	//     logger.Log.Debugf("Type=Response, ID=%d, UID=%s, MID=%d, Data=%dbytes",
	//          0, 0, mid, len(d))
	// default:
	//     logger.Log.Infof("Type=Response, ID=%d, UID=%s, MID=%d, Data=%+v",
	//         0, 0, mid, v)
	// }

	return s.send(pendingMessage{typ: message.Response, mid: mid, payload: v, err: err})
}

// ---- write over ----
// ---- read begin ----
func (s *ClientSession) read() {
	conn := s.conn
	for {
		if s.GetStatus() == StatusClosed {
			break
		}
		msg, err := conn.GetNextMessage()

		if err != nil {
			if err != constants.ErrConnectionClosed {
				logger.Log.Errorf("Error reading next available message: %s", err.Error())
			}
			s.Close()
			return
		}

		packets, err := s.cfg.Decoder.Decode(msg)
		if err != nil {
			logger.Log.Errorf("Failed to decode message: %s", err.Error())
			return
		}

		if len(packets) < 1 {
			logger.Log.Warnf("Read no packets, data: %v", msg)
			continue
		}

		// process all packet
		for i := range packets {
			if err := s.processPacket(packets[i]); err != nil {
				logger.Log.Errorf("Failed to process packet: %s", err.Error())
				return
			}
		}
	}
}

func (s *ClientSession) processPacket(p *packet.Packet) error {
	switch p.Type {
	case packet.Handshake:
		logger.Log.Debug("Received handshake packet")
		if err := s.SendHandshakeResponse(); err != nil {
			logger.Log.Errorf("Error sending handshake response: %s", err.Error())
			return err
		}
		//logger.Log.Debugf("Session handshake Id=%d, Remote=%s", a.GetSession().ID(), a.RemoteAddr())

		// Parse the json sent with the handshake by the client
		handshakeData := &HandshakeData{}
		err := json.Unmarshal(p.Data, handshakeData)
		if err != nil {
			s.SetStatus(StatusClosed)
			return fmt.Errorf("Invalid handshake data. Id=%d", 0)
		}

		//a.GetSession().SetHandshakeData(handshakeData)
		s.SetStatus(StatusHandshake)
		//err = a.GetSession().Set(constants.IPVersionKey, a.IPVersion())
		// if err != nil {
		//     logger.Log.Warnf("failed to save ip version on session: %q\n", err)
		// }

		logger.Log.Debug("Successfully saved handshake data")

	case packet.HandshakeAck:
		s.lastHeartBeat = common.NowMs()
		s.SetStatus(StatusWorking)
		//logger.Log.Debugf("Receive handshake ACK Id=%d, Remote=%s", a.GetSession().ID(), a.RemoteAddr())

	case packet.Data:
		if s.GetStatus() < StatusWorking {
			// return fmt.Errorf("receive data on socket which is not yet ACK, session will be closed immediately, remote=%s",
			//     a.RemoteAddr().String())
			return nil
		}

		msg, err := message.Decode(p.Data)
		if err != nil {
			return err
		}
		s.processMessage(msg)
		//logger.Log.Debugf("%v", msg)

	// 心跳处理
	// 主要是断开不活跃的连接
	case packet.Heartbeat:

		s.lastHeartBeat = common.NowMs()
		//logger.Log.Infof(" %v got heartbeat %v", s.netId, s.lastHeartBeat)
	}

	//a.SetLastAt()
	return nil
}

// ---- read over ----
// ---- heartbeat ----
func (s *ClientSession) heartbeat() {
	// timer
	t := time.NewTicker(DefaultHeartbeatTime)
	defer func() {
		t.Stop()
	}()

	for {
		select {
		case <-s.chanClose:
			return
		case <-t.C:
			s.checkHeartBeatTimeout()
			s.sendHeartBeat()
		}
	}
}

func (s *ClientSession) checkHeartBeatTimeout() {
	if s.GetStatus() != StatusWorking {
		return
	}

	if common.NowMs() < s.lastHeartBeat+2*DefaultHeartbeatTimeSeconds*1000 {
		return
	}

	logger.Log.Infof(" %v heartbeat timeout and close %v >= %v", s.netId, common.NowMs(), s.lastHeartBeat)
	// close s
	s.Close()
}

func (s *ClientSession) sendHeartBeat() {
	if s.GetStatus() != StatusWorking {
		return
	}
	pWrite := &pendingWrite{
		data: hbd,
	}

	s.pushToSend(pWrite)
}

func (s *ClientSession) pushToSend(p *pendingWrite) {
	defer func() {
		// maybe closed
		if e := recover(); e != nil {
		}
	}()
	s.chSend <- p
}

// ---- heartbeat over ----

// TODO: 流量控制
func (s *ClientSession) processMessage(msg *message.Message) {
	s.GetImpl().ProcessMessage(s, msg)
}

// 生成一个hdd packet
func hbdEncode(heartbeatTimeout time.Duration, packetEncoder codec.PacketEncoder, dataCompression bool, serializerName string) {
	hData := map[string]interface{}{
		"code": 200,
		"sys": map[string]interface{}{
			"heartbeat":  heartbeatTimeout.Seconds(),
			"dict":       message.GetDictionary(),
			"serializer": serializerName,
		},
	}
	data, err := json.Marshal(hData)
	if err != nil {
		panic(err)
	}

	if dataCompression {
		logger.Log.Debugf("hbdEncode compression")
		compressedData, err := compression.DeflateData(data)
		if err != nil {
			panic(err)
		}

		if len(compressedData) < len(data) {
			data = compressedData
		}
	}

	hrd, err = packetEncoder.Encode(packet.Handshake, data)
	if err != nil {
		panic(err)
	}

	hbd, err = packetEncoder.Encode(packet.Heartbeat, nil)
	if err != nil {
		panic(err)
	}
}
