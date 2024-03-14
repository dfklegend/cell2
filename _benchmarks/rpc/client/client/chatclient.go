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
	"github.com/dfklegend/cell2/utils/profiler"
	"github.com/dfklegend/cell2/utils/runservice"
	json2 "github.com/dfklegend/cell2/utils/serialize/json"
	"github.com/dfklegend/cell2/utils/timer"

	"client/protos"
)

const (
	STATE_INIT = iota
	STATE_CONNECT_GATE
	STATE_TESTING
	STATE_FINISH
)

var (
	// speed msgs / milisecond
	SendSpeed  = float64(120)
	MaxRequest = 10000

	ErrorLogInterval int64 = 0
)

type ChatClient struct {
	ID         string
	Client     *client.CellClient
	RunService *runservice.StandardRunService
	state      int
	msgBind    bool

	beginTime    int64
	lastSend     int64
	sendOverTime int64
	requestSent  int
	responseGot  int
	errorGot     int
	// 检查读取顺序
	prevNumber int

	// stat
	totalCost int64

	maxRequest int
	timerDo    timer.IdType

	stat         *profiler.ReqStat
	sendTimeRest float64

	nextLogSendError int64
	detailLog        bool
}

func NewChatClient(name string) *ChatClient {
	c := client.NewCellClient(name, json2.GetDefaultSerializer())
	r := c.GetRunService()

	return &ChatClient{
		ID:         name,
		Client:     c,
		RunService: r,
		state:      STATE_INIT,
		msgBind:    false,
		stat:       profiler.NewReqStat(),
		detailLog:  false,
	}
}

func (cc *ChatClient) SetDetailLog(v bool) {
	cc.detailLog = v
	cc.Client.SetDetailLog(v)
}

func (cc *ChatClient) setState(s int) {
	cc.state = s
}

func (cc *ChatClient) getState() int {
	return cc.state
}

func (cc *ChatClient) Start(address string) {
	c := cc.Client
	c.Start(address)
	c.SetRetryDelay(5)

	c.SetCBConnected(cc.onConnected)
	c.SetCBBreak(func() {
		logger.Log.Errorf("connection break")
		cc.setState(STATE_INIT)
		cc.closeTestTimer()
	})
	cc.bindMsgs()
}

func (cc *ChatClient) Stop() {
	cc.Client.Stop()
}

func (cc *ChatClient) onConnected() {
	logger.Log.Debugf("onConnected")
	if cc.getState() != STATE_INIT {
		return
	}

	// 开始
	cc.queryGate()
	cc.setState(STATE_CONNECT_GATE)
}

func (cc *ChatClient) bindMsgs() {
}

func (cc *ChatClient) queryGate() {
	logger.Log.Debugf("queryGate")

	c := cc.Client
	c.GetClient().SendRequest("gate.gate.querygate", []byte("{}"), func(error bool, msg *message.Message) {
		fmt.Println("ack from cb")
		fmt.Println(string(msg.Data))

		if error {
			fmt.Println("query gate error ")
			return
		}

		// TODO:解析，使用正确端口
		port := cc.getGatePort(msg.Data)
		c.Start(fmt.Sprintf("127.0.0.1:%v", port))
		c.WaitReady()

		m := make(map[string]interface{})
		m["name"] = "haha"
		c.GetClient().SendRequest("gate.gate.login", []byte(common.SafeJsonMarshal(m)), cc.onLoginRet)
	})
}

func (cc *ChatClient) getGatePort(data []byte) string {
	obj := protos.QueryGateAck{}
	json.Unmarshal(data, &obj)

	port := obj.Port
	subs := strings.Split(port, ",")
	if len(subs) >= 2 {
		return subs[1]
	}
	return ""
}

func (cc *ChatClient) onLoginRet(error bool, msg *message.Message) {
	if error {
		logger.Log.Errorf("login ret error")

		// delay restart
		cc.Client.GetRunService().GetTimerMgr().After(3*time.Second, func(args ...any) {
			cc.Client.Disconnect()
		})
		return
	}
	// 可以发送消息
	obj := protos.NormalAck{}
	json.Unmarshal(msg.Data, &obj)
	logger.Log.Debugf("onLoginRet: %+v", obj)

	//c := cc.Client
	//m := make(map[string]interface{})
	//m["name"] = "haha"
	//m["content"] = "hello"
	//data := []byte(common.SafeJsonMarshal(m))
	//c.GetClient().SendRequest("chat.chat.sendchat", data, func(error bool, msg *message.Message) {
	//	logger.Log.Infof("send chat ret, error:%v", error)
	//})

	cc.startTest()
}

func (cc *ChatClient) startTest() {
	c := cc.Client

	cc.setState(STATE_TESTING)

	cc.beginTime = common.NowMs()
	cc.lastSend = cc.beginTime
	cc.requestSent = 0
	cc.responseGot = 0
	cc.errorGot = 0
	cc.prevNumber = 0
	cc.maxRequest = MaxRequest
	cc.sendTimeRest = 0

	cc.stat.Begin()

	cc.timerDo = c.GetRunService().GetTimerMgr().AddTimer(time.Millisecond, func(args ...any) {
		cc.doTest()
	})
}

func makeData(index int, log int) []byte {
	return []byte(common.SafeJsonMarshal(&protos.Hello{
		Msg:    fmt.Sprintf("Hello from %v", index),
		Number: index,
		Log:    log,
	}))
}

func (cc *ChatClient) closeTestTimer() {
	if cc.timerDo == 0 {
		return
	}
	cc.Client.GetRunService().GetTimerMgr().Cancel(cc.timerDo)
	cc.timerDo = 0
}

func (cc *ChatClient) doTest() {
	// 中间断了
	if !cc.Client.IsReady() {
		logger.Log.Infof("---- doTest break %v", cc.ID)
		cc.closeTestTimer()
		return
	}

	if cc.IsWaitResponseTimeout() {
		cc.closeTestTimer()
		logger.Log.Infof("---- timeout will restart later %v", cc.ID)
		cc.Client.GetRunService().GetTimerMgr().After(3*time.Second, func(args ...any) {
			cc.Client.Disconnect()
		})
		return
	}

	if cc.IsAllFinish() {
		logger.Log.Infof("finished!")
		cc.setState(STATE_FINISH)
		cc.closeTestTimer()
		cc.stat.End()
		cc.stat.Report(cc.ID)
		cc.calcStat()
		cc.dumpStat()
		// dump stat

		logger.Log.Infof("---- will restart later %v", cc.ID)
		cc.Client.GetRunService().GetTimerMgr().After(3*time.Second, func(args ...any) {
			cc.Client.Disconnect()
		})
		return
	}

	now := common.NowMs()
	if now <= cc.lastSend {
		return
	}

	times := float64(now-cc.lastSend)*SendSpeed + cc.sendTimeRest

	// flow ctrl
	if times > SendSpeed*10 && times > 100 {
		times = SendSpeed
	}

	cc.lastSend = now

	if times < 1 {
		cc.sendTimeRest = times
		return
	}
	// 余量
	cc.sendTimeRest = 0

	for i := int64(0); i < int64(times) && cc.requestSent < cc.maxRequest; i++ {
		cc.doSend(cc.requestSent)
	}
}

func (cc *ChatClient) doSend(index int) {

	beginTime := common.NowMs()
	logEnable := 0
	if cc.detailLog {
		logEnable = 1
	}
	cc.Client.GetClient().SendRequest("chat.chat.hello",
		makeData(index, logEnable), func(error bool, msg *message.Message) {
			if error {
				cc.errorGot++
				// 避免出错了，太多输出
				if true || common.NowMs() >= cc.nextLogSendError {
					logger.Log.Errorf("%v error SendRequest errorGot: %v", cc.ID, cc.errorGot)
					cc.nextLogSendError = common.NowMs() + ErrorLogInterval
				}
				return
			}
			cc.responseGot++
			cc.stat.AddCall(common.NowMs() - beginTime)

			obj := protos.NormalAck{}
			json.Unmarshal(msg.Data, &obj)

			responseGot := obj.Code
			if cc.detailLog || responseGot%1000 == 0 {
				logger.Log.Infof("%v responseGot : %v", cc.ID, responseGot)
			}

			if obj.Code != index {
				logger.Log.Errorf("%v MISMATCHED!!!", cc.ID)
			}

			// 如果前面发生超时，后面肯定MisOrder
			// 检查读取顺序
			if cc.errorGot == 0 {
				if /*cc.detailLog &&*/ responseGot != cc.prevNumber+1 && cc.prevNumber > 0 {
					if true || common.NowMs() >= cc.nextLogSendError {
						logger.Log.Errorf("%v MISORDER!!!, %v(want) != %v(actual)", cc.ID, cc.prevNumber+1, responseGot)
						cc.nextLogSendError = common.NowMs() + ErrorLogInterval
					}
				}
				cc.prevNumber = obj.Code
			}

		})
	cc.requestSent++
	if cc.requestSent == cc.maxRequest {
		cc.sendOverTime = common.NowMs()
	}

	if cc.detailLog || index%1000 == 0 {
		logger.Log.Infof("%v requestSent : %v", cc.ID, index)
	}
}

func (cc *ChatClient) IsAllFinish() bool {
	return cc.getState() == STATE_TESTING &&
		cc.requestSent == cc.maxRequest && cc.requestSent == cc.responseGot
}

func (cc *ChatClient) IsWaitResponseTimeout() bool {
	return cc.getState() == STATE_TESTING &&
		cc.requestSent == cc.maxRequest && common.NowMs() > cc.sendOverTime+30*1000
}

func (cc *ChatClient) calcStat() {
	cc.totalCost = common.NowMs() - cc.beginTime
}

func (cc *ChatClient) dumpStat() {
	if cc.responseGot == 0 {
		return
	}
	l := logger.Log
	l.Infof(" avg: %v ms/req  %v req/second", float64(cc.totalCost)/float64(cc.responseGot),
		int64(cc.responseGot*1000)/cc.totalCost)

}
