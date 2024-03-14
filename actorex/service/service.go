// 扩展Service
// 提供Request语法
// 提供将调用映射到具体处理者的功能

package service

import (
	"errors"
	"log"
	"time"

	"github.com/dfklegend/cell2/utils/common"
	"github.com/dfklegend/cell2/utils/runservice"
	"github.com/dfklegend/cell2/utils/timer"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"

	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	l "github.com/dfklegend/cell2/utils/logger"
)

const (
	DefaultSerializeId       = 0 // DefaultSerializeId [protoactor]/remote/serializer
	MaxReqId           int32 = 0x7FFFFFF0
	RequestTimeout           = 30 * 1000
)

var (
	ErrTimeout = errors.New("request timeout")
	showTrace  = false
)

var (
	// 测试关闭发送，测试超时
	testReqTimeoutEnable = false
)

type RequestWaitResponse struct {
	ReqId   int32
	Route   string
	Msg     interface{}
	Timeout int64
	CB      ResCBFunc
}

// 	Service可以发起Request
// 	当Request返回时，会调用cb
// 	Request的结构体，必须能序列化
//	TODO:	优化,根据是否远程对象，来确定是否要序列化
type Service struct {
	extProps      *ExtProps
	reqReceiver   IRequestReceiver
	apiDispatcher IAPIDispatcher
	// 简单存一下cb
	Handlers map[int32]*RequestWaitResponse
	nextId   int32
	Context  actor.Context
	// 必须搭配StandardRunService
	runService *runservice.StandardRunService

	timerCheckExpired timer.IdType
}

func NewService() *Service {
	s := &Service{
		Handlers: make(map[int32]*RequestWaitResponse),
	}
	s.InitReqReceiver(s)
	return s
}

func (s *Service) InitReqReceiver(receiver IRequestReceiver) {
	s.reqReceiver = receiver
}

func (s *Service) SetAPIDispatcher(disp IAPIDispatcher) {
	s.apiDispatcher = disp
}

// 	IService

func (s *Service) SetExtProps(props *ExtProps) {
	s.extProps = props
}

func (s *Service) SetRunService(rservice *runservice.StandardRunService) {
	s.runService = rservice
}

func (s *Service) GetRunService() *runservice.StandardRunService {
	return s.runService
}

func (s *Service) AllocReqId() int32 {
	if s.nextId >= MaxReqId {
		s.nextId = 0
	}

	s.nextId += 1
	return s.nextId
}

func (s *Service) OnCreate() {
}

func (s *Service) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		s.onStarted(ctx)
	case *actor.Stop:
		s.onStop()
	case *messages.ServiceRequest:
		s.handleRequest(ctx, msg)
	case *messages.ServiceResponse:
		s.handleResponse(msg)
	}
}

func (s *Service) onStarted(ctx actor.Context) {
	s.Context = ctx
	if s.extProps != nil {
		s.extProps.doPostStartFuncs(ctx)
	}
}

func (s *Service) onStop() {
	// 不考虑共享runService
	s.runService.Stop()
}

func (s *Service) calcTimeout() int64 {
	if testReqTimeoutEnable {
		return common.NowMs() + 3*1000
	}
	return common.NowMs() + RequestTimeout
}

//	调用
// 	线程安全性：
//	必须在service环境内调用request
//	基于actor模式，不希望service的request在非本service执行环境内调用
func (s *Service) Request(pid *actor.PID, msg interface{}, cbFunc ResCBFunc) {
	s.doRequestEx(pid, "", msg, true, cbFunc)
}

func (s *Service) RequestEx(pid *actor.PID, route string, msg interface{}, cbFunc ResCBFunc) {
	s.doRequestEx(pid, route, msg, true, cbFunc)
}

func (s *Service) Notify(pid *actor.PID, msg interface{}) {
	s.doRequestEx(pid, "", msg, false, nil)
}

func (s *Service) NotifyEx(pid *actor.PID, route string, msg interface{}) {
	s.doRequestEx(pid, route, msg, false, nil)
}

//	Notify就是无需返回值的Request
func (s *Service) doRequestEx(pid *actor.PID, route string, msg interface{}, isRequest bool, cbFunc ResCBFunc) {
	//. 分配reqId
	//. 记录cb
	//. 调用
	var reqId int32
	reqId = NotifyReqID

	if isRequest {
		reqId = s.AllocReqId()

		wait := &RequestWaitResponse{
			ReqId:   reqId,
			Route:   route,
			Msg:     msg,
			CB:      cbFunc,
			Timeout: s.calcTimeout(),
		}

		s.Handlers[reqId] = wait
	}

	bytes, typeName, err := remote.Serialize(msg, DefaultSerializeId)
	if err != nil {
		l.Log.Errorf("%v msg serialize failed: %v", route, err)
		return
	}

	req := &messages.ServiceRequest{
		Sender: s.Context.Self(),
		ReqId:  reqId,
		Route:  route,
		Type:   typeName,
		Body:   bytes,
	}

	s.sendRequest(pid, req)

	if showTrace {
		if isRequest {
			l.Log.Infof("send Request:%v \n", req.ReqId)
		}
	}
	// 根据情况看是否启动检查Timer
	if isRequest {
		s.tryStartCheckTimer()
	}
}

func (s *Service) sendRequest(pid *actor.PID, req *messages.ServiceRequest) {
	if testReqTimeoutEnable {
		l.Log.Infof("this request: %v send disabled for test\n", req.ReqId)
		return
	}
	s.Context.Send(pid, req)
}

func (s *Service) tryStartCheckTimer() {
	if s.timerCheckExpired > 0 {
		return
	}

	s.timerCheckExpired = s.runService.GetTimerMgr().AddTimer(
		time.Second, func(args ...interface{}) {
			s.checkExpired()
		})
}

func (s *Service) freeTimer() {
	s.runService.GetTimerMgr().Cancel(s.timerCheckExpired)
	s.timerCheckExpired = 0
}

func (s *Service) checkExpired() {
	if len(s.Handlers) == 0 {
		s.freeTimer()
		return
	}

	now := common.NowMs()
	expires := make([]int32, 0)
	for id, one := range s.Handlers {
		if one.Timeout < now {
			expires = append(expires, id)
		}
	}

	if len(expires) == 0 {
		return
	}

	for _, reqId := range expires {
		wait := s.Handlers[reqId]
		l.Log.Infof("request timeout: reqId: %v route: %v argType: %T, arg: {%v}\n",
			reqId, wait.Route, wait.Msg, wait.Msg)

		if wait.CB != nil {
			wait.CB(ErrTimeout, nil)
		}
		delete(s.Handlers, reqId)
	}

}

//	Response 返回信息
//  Response, ResponseEx线程安全
func (s *Service) Response(req *messages.ServiceRequest, errCode int32, errInfo string, msg interface{}) {
	s.ResponseEx(req, errCode, errInfo, msg, true)
}

func (s *Service) ResponseEx(req *messages.ServiceRequest, errCode int32, errInfo string,
	msg interface{}, needSerialze bool) {

	// need no response
	if req.ReqId == NotifyReqID {
		return
	}

	if req.Sender == nil {
		return
	}

	res := &messages.ServiceResponse{
		ReqId: req.ReqId,
	}

	if errCode != 0 {
		res.ErrCode = errCode
		res.ErrInfo = errInfo
	} else {
		// 有返回值
		if needSerialze && msg != nil {
			bytes, typeName, err := remote.Serialize(msg, DefaultSerializeId)
			if err != nil {
				panic(err)
			}
			res.Type = typeName
			res.Body = bytes
		}
	}

	s.Context.Send(req.Sender, res)
}

func (s *Service) handleRequest(ctx actor.Context, request *messages.ServiceRequest) {
	// 检查是否通过apimapper来dispatcher
	if request.Route != "" {
		if s.apiDispatcher != nil &&
			s.apiDispatcher.Dispatch(ctx, s, request.Route, request) {
			// 被处理了
			return
		}
	}

	msg, err := remote.Deserialize(request.Body, request.Type, DefaultSerializeId)
	if err != nil {
		panic(err)
	}

	if s.reqReceiver == nil {
		return
	}
	s.reqReceiver.ReceiveRequest(ctx, request, msg)
}

// 	to override ReceiveRequest
func (s *Service) ReceiveRequest(ctx actor.Context, request *messages.ServiceRequest, rawMsg interface{}) {

}

// 接收到response
func (s *Service) handleResponse(response *messages.ServiceResponse) {
	if showTrace {
		log.Printf("got response:%v \n", response.ReqId)
	}
	wait := s.Handlers[response.ReqId]
	if wait == nil {
		log.Printf("miss response: %v\n", response.ReqId)
		return
	}

	var err error
	var msg any
	if response.ErrCode != CodeSucc {
		err = errors.New(response.ErrInfo)
	} else {
		if response.Type != "" {
			msg, err = remote.Deserialize(response.Body, response.Type, DefaultSerializeId)
		}
	}

	if wait.CB != nil {
		wait.CB(err, msg)
	}

	// remove it
	delete(s.Handlers, response.ReqId)
}

func (s *Service) Send(pid *actor.PID, msg interface{}) {
	s.Context.Send(pid, msg)
}

// Post 比如，可以post到service环境内request
func (s *Service) Post(cb func()) {
	s.GetRunService().GetScheduler().Post(cb)
}
