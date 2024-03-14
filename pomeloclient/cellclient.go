package client

import (
	"errors"
	"reflect"
	"time"

	"github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/apimapper/apientry"
	"github.com/dfklegend/cell2/pomelonet/common/conn/message"

	nclient "github.com/dfklegend/cell2/pomelonet/client"
	"github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/runservice"
	"github.com/dfklegend/cell2/utils/sche"
	"github.com/dfklegend/cell2/utils/serialize"
	"github.com/dfklegend/cell2/utils/serialize/json"
)

const (
	STATE_INIT = iota
	STATE_CONNECTING
	STATE_CONNECTED
	STATE_BROKEN
	STATE_PENDINGRETRY // 等待5s再重连
)

var (
	ErrorRequestError = errors.New("Request error")
)

type HandleFunc func()

/*
	CellClient
	go client实现
	可以发送请求，可以注册消息处理函数
*/
type CellClient struct {
	name       string
	TheClient  *nclient.Client
	runService *runservice.StandardRunService
	state      int
	firstStart bool
	autoRetry  bool

	connectingTime int
	retryWait      int
	retryDelay     int
	tarAddress     string

	cbBreak     HandleFunc
	cbConnected HandleFunc
	detailLog   bool

	collection *apientry.APICollection

	serializer serialize.Serializer
	owner      IClientOwner
}

func NewCellClient(name string, serializer serialize.Serializer) *CellClient {
	if serializer == nil {
		serializer = json.GetDefaultSerializer()
	}
	c := &CellClient{
		name:       name,
		TheClient:  nclient.New(),
		runService: runservice.NewStandardRunService(name),
		firstStart: true,
		autoRetry:  true,
		state:      STATE_INIT,
		retryDelay: 5,
		detailLog:  false,
		collection: apientry.NewCollection(),
		serializer: serializer,
	}
	c.TheClient.SetName(name)
	return c
}

func (c *CellClient) Reserve() {
}

func (c *CellClient) Handle() {
}

func (c *CellClient) GetOwner() IClientOwner {
	return c.owner
}

func (c *CellClient) SetOwner(owner IClientOwner) {
	c.owner = owner
}

func (c *CellClient) SetDetailLog(v bool) {
	c.detailLog = v
	c.TheClient.SetDetailLog(v)
}

func (c *CellClient) GetRunService() *runservice.StandardRunService {
	return c.runService
}

func (c *CellClient) SetRetryDelay(delay int) {
	c.retryDelay = delay
}

func (c *CellClient) setState(state int) {
	c.state = state
}

func (c *CellClient) getState() int {
	return c.state
}

func (c *CellClient) IsReady() bool {
	return c.getState() == STATE_CONNECTED
}

func (c *CellClient) SetCBBreak(cb HandleFunc) {
	c.cbBreak = cb
}

func (c *CellClient) SetCBConnected(cb HandleFunc) {
	c.cbConnected = cb
}

func (c *CellClient) Start(address string) {
	if c.firstStart {
		c.runService.Start()
		c.runService.GetEventCenter().SetLocalUseChan(true)
		c.addUpdate()

		c.firstStart = false
	}
	c.setState(STATE_INIT)
	c.autoRetry = true

	c.tarAddress = address
	c.Connect(address)
}

func (c *CellClient) Connect(address string) {
	old := c.TheClient
	if old != nil {
		old.Disconnect()
	}
	c.TheClient = nclient.New()
	c.TheClient.SetName(c.name)
	c.TheClient.SetDetailLog(c.detailLog)

	c.setState(STATE_CONNECTING)
	c.TheClient.ConnectTo(address)

	// TODO: 移除掉老的MsgChannel selector
	c.addMsgProcess()
}

func (c *CellClient) WaitReady() {
	for !c.TheClient.Ready {
		time.Sleep(100 * time.Millisecond)
		c.onUpdate()
	}
}

func (c *CellClient) Disconnect() {
	c.TheClient.Disconnect()
}

func (c *CellClient) Stop() {
	c.autoRetry = false
	c.TheClient.Disconnect()
	c.runService.Stop()
}

func (c *CellClient) GetClient() *nclient.Client {
	return c.TheClient
}

// SendRequest 必须提供返回值参数类型，不然无法序列化返回值
// 相当于由调用者决定应该是什么返回值，如果不对应，就是错误结果
func (c *CellClient) SendRequest(route string, msg any, cb apientry.HandlerCBFunc, retType reflect.Type) (uint, error) {
	if retType != nil && retType.Kind() != reflect.Ptr {
		logger.Log.Debugf("retType 必须是指针或者Nil")
		return 0, nil
	}

	data, err := c.serializer.Marshal(msg)
	if err != nil {
		logger.Log.Debugf("Marshal error: %v", err)
		return 0, err
	}

	requestId, err1 := c.GetClient().SendRequest(route, data, func(error bool, msg *message.Message) {
		if error {
			apientry.CheckInvokeCBFunc(cb, ErrorRequestError, nil)
			return
		}
		// TODO: 检查retType是否指针
		if retType != nil {
			ret := reflect.New(retType.Elem()).Interface()
			c.serializer.Unmarshal(msg.Data, ret)
			apientry.CheckInvokeCBFunc(cb, nil, ret)
		} else {
			apientry.CheckInvokeCBFunc(cb, nil, nil)
		}

	})
	return requestId, err1
}

func (c *CellClient) SendNotify(route string, msg any) error {
	data, err := c.serializer.Marshal(msg)
	if err != nil {
		logger.Log.Debugf("Marshal error: %v", err)
		return err
	}

	return c.GetClient().SendNotify(route, data)
}

func (c *CellClient) RegisterDefaultEntry(e apimapper.IAPIEntry, options ...apientry.Option) {
	newOptions := append([]apientry.Option{}, options...)
	newOptions = append(newOptions, apientry.WithInnerGroupName())
	c.collection.Register(e, newOptions...)
}

func (c *CellClient) RegisterEntry(e apimapper.IAPIEntry, options ...apientry.Option) {
	c.collection.Register(e, options...)
}

func (c *CellClient) BuildEntries() {
	c.collection.Build()
}

func (c *CellClient) addMsgProcess() {
	selector := c.runService.GetSelector()
	selector.AddSelector("cellclient", sche.NewFuncSelector(reflect.ValueOf(c.TheClient.MsgChannel()),
		func(v reflect.Value, recvOk bool) {
			if !recvOk {
				return
			}

			msg := v.Interface().(*nclient.ClientMsg)
			c.processMsg(msg)
		}))
}

func (c *CellClient) processMsg(msg *nclient.ClientMsg) {
	if msg.Cb != nil {
		c.processResponse(msg)
	} else {
		// push msg
		logger.Log.Debugf("got push:%v", msg.Msg)
		//Msg := msg.Msg
		//c.runService.GetEventCenter().Publish(Msg.Route, Msg.Data)
		c.processPushMsg(msg)
	}
}

func (c *CellClient) processResponse(msg *nclient.ClientMsg) {
	msg.Cb(msg.Msg.Err, msg.Msg)
}

func (c *CellClient) processPushMsg(msg *nclient.ClientMsg) {
	logger.Log.Debugf("got push:%v", msg.Msg)
	apientry.CallWithSerialize(c.collection, c, msg.Msg.Route, msg.Msg.Data, nil, c.serializer)
}

func (c *CellClient) addUpdate() {
	c.runService.GetTimerMgr().AddTimer(1*time.Second, func(args ...interface{}) {
		c.onUpdate()
	})
}

func (c *CellClient) onUpdate() {
	switch c.getState() {
	case STATE_CONNECTING:
		c.onConnecting()
	case STATE_CONNECTED:
		if !c.TheClient.Connected {
			c.setState(STATE_BROKEN)
			c.onBreak()
		}
	case STATE_PENDINGRETRY:
		c.onPendingRetry()
	}
}

func (c *CellClient) onConnecting() {
	if c.TheClient.Ready {
		c.setState(STATE_CONNECTED)
		if c.cbConnected != nil {
			c.cbConnected()
		}
		return
	}
	c.connectingTime++
	if c.connectingTime > 5 {
		c.onBreak()
	}
}

func (c *CellClient) onBreak() {
	if c.cbBreak != nil {
		c.cbBreak()
	}
	// 看要不要重连
	if c.autoRetry {
		c.beginRetry()
	}
}

func (c *CellClient) beginRetry() {
	c.setState(STATE_PENDINGRETRY)
	c.retryWait = c.retryDelay
}

func (c *CellClient) onPendingRetry() {
	c.retryWait--
	logger.Log.Debugf("reconnect wait:%v", c.retryWait)
	if c.retryWait <= 0 {
		logger.Log.Debugf("reconnect")
		c.Connect(c.tarAddress)
	}
}
