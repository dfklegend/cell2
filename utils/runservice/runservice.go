package runservice

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"time"

	"github.com/dfklegend/cell2/utils/common"
	"github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/sche"
)

// RunService的
var (
	idService *common.SerialIdService
	rsScheMgr *sche.Mgr

	PerfLogLevel int32 = LevelNormal
)

const (
	LevelDisable int32 = iota
	LevelNormal        // 输出超过了25ms的帧
	LevelDetail        // 详细输出信息，调试用
)

const (
	OneStepWarnning int64 = 25 * 1000000
	CostWarnning0   int64 = 30
	CostWarnning1   int64 = 50
	CostWarnning2   int64 = 100

	OutputPutInterval int64 = 10
)

func init() {
	rsScheMgr = sche.DefaultScheMgr
	idService = common.NewSerialIdService()
}

func SetPerfLogLevel(level int32) {
	PerfLogLevel = level
}

func GetScheMgr() *sche.Mgr {
	return rsScheMgr
}

// 	代表一个可执行的service
// 	每个一个单独的routine
// 	是一个可投递目标
// 	可以用selector来定制循环
type RunService struct {
	Name string
	// 调度器
	scheduler *sche.Sche
	// 用于添加select队列
	selector *sche.MultiSelector
	running  bool

	chanClose chan int

	// running state
	busyTime int64

	nextOutputInfo  int64
	totalSkipCost   int64
	skipOutputTimes int

	// 本地 blackboard
	// 可以集成一些对象
	vars map[string]any
}

func makeName(name string) string {
	if name != "" {
		return name
	}
	return fmt.Sprintf("rs%v", idService.AllocId())
}

func NewRunService(name string) *RunService {
	name = makeName(name)
	return &RunService{
		Name:      name,
		scheduler: rsScheMgr.GetSche(name),
		selector:  sche.NewMultiSelector(),
		running:   true,
		chanClose: make(chan int, 1),
	}
}

func (r *RunService) GetScheduler() *sche.Sche {
	return r.scheduler
}

func (r *RunService) GetSelector() *sche.MultiSelector {
	return r.selector
}

func (r *RunService) Start() {
	r.addSchedulerSelector()
	r.addCloseChan()

	go r.loop()
}

func (r *RunService) Stop() {
	// 通知跳出loop
	close(r.chanClose)
	r.scheduler.Stop()
	r.selector.Stop()
	GetScheMgr().DelSche(makeName(r.Name))
}

func (r *RunService) IsStopped() bool {
	return !r.running
}

func (r *RunService) addSchedulerSelector() {
	scheduler := r.scheduler
	selector := r.selector
	selector.AddSelector("sheduler", sche.NewFuncSelector(reflect.ValueOf(scheduler.GetChanTask()),
		func(v reflect.Value, recvOk bool) {
			if !recvOk {
				return
			}

			task := v.Interface().(*sche.RunTask)
			scheduler.DoTask(task)
		}))
}

func (r *RunService) addCloseChan() {
	selector := r.selector
	selector.AddSelector("__close__", sche.NewFuncSelector(reflect.ValueOf(r.chanClose),
		func(v reflect.Value, recvOk bool) {
			r.running = false
		}))
}

func (r *RunService) loop() {
	defer func() {
		log.Println("RunServeice loop end")
	}()

	for r.running {
		r.selector.HandleOnce()
		r.analysisRunning()
	}
}

// 统计方式
// 传统计时可能由于 runtime.Gosched() 不准确
func (r *RunService) analysisRunning() {
	s := r.selector

	if PerfLogLevel >= LevelDetail {
		logger.Log.Infof("runservice :%v lastcost  %v do cost: %v ns select cost: %v ns",
			r.Name, s.LastDoChannelName, s.LastDoCost, s.LastSelectCost)
	}

	//
	if PerfLogLevel >= LevelNormal && s.LastDoCost > OneStepWarnning {
	}

	// 等待过(释放了CPU)
	selectWaited := s.LastSelectCost > 1000
	if selectWaited {
		r.busyTime = s.LastDoCost
	} else {
		r.busyTime += s.LastDoCost
	}

	// 30ms
	cost := r.busyTime / 1000000
	if cost >= CostWarnning0 {
		if PerfLogLevel >= LevelNormal {
			// 避免太频繁
			// 3s输出一次
			nowMs := common.NowMs()
			if nowMs > r.nextOutputInfo {
				logger.Log.Infof("runservice :%v  heavy frame <channel:%v> do cost: %v ms ", r.Name, s.LastDoChannelName, s.LastDoCost/1000000)
				//logger.Log.Warnf("runservice :%v is busy, cost: %v", r.Name, cost)
				r.nextOutputInfo = nowMs + OutputPutInterval*1000
				if r.skipOutputTimes > 0 {
					logger.Log.Infof("runservice :%v  skip output heavy frame times: %v in %v seconds avg cost: %v ms",
						r.Name, r.skipOutputTimes, OutputPutInterval, r.totalSkipCost/(int64)(r.skipOutputTimes))
					r.skipOutputTimes = 0
					r.totalSkipCost = 0
				}
			} else {
				r.skipOutputTimes++
				r.totalSkipCost += s.LastDoCost / 1000000
			}
		}
		r.busyTime = 0
	}

	// 调度
	if cost >= CostWarnning2 {
		time.Sleep(2 * time.Millisecond)
	} else if cost >= CostWarnning1 {
		time.Sleep(1 * time.Millisecond)
	} else if cost >= CostWarnning0 {
		runtime.Gosched()
	}
}

func (r *RunService) SetValue(key string, v any) {
	r.vars[key] = v
}

func (r *RunService) GetValue(key string) any {
	return r.vars[key]
}
