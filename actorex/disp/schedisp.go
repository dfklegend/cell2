package disp

import (
	"reflect"

	"github.com/dfklegend/cell2/utils/runservice"
	"github.com/dfklegend/cell2/utils/sche"
)

// 使用runservice来进行调度
// 保证在稳定的携程调度
// implement actor.Dispatcher
type scheDisp struct {
	runService *runservice.StandardRunService
	chanTask   ChanTask
}

func NewScheDisp(name string) *scheDisp {
	return &scheDisp{
		runService: runservice.NewStandardRunService(name),
		chanTask:   make(ChanTask, 9),
	}
}

func (s *scheDisp) GetRunService() *runservice.StandardRunService {
	return s.runService
}

func (s *scheDisp) Schedule(fn func()) {
	s.chanTask <- fn
}

func (s *scheDisp) Throughput() int {
	return 99
}

func (s *scheDisp) Start() {
	s.addFunSelector()
	s.runService.Start()
}

func (s *scheDisp) Stop() {
	s.runService.Stop()
}

func (s *scheDisp) addFunSelector() {
	selector := s.runService.GetSelector()
	selector.AddSelector("scheDisp", sche.NewFuncSelector(reflect.ValueOf(s.chanTask),
		func(v reflect.Value, recvOk bool) {
			if !recvOk {
				return
			}

			fn := v.Interface().(func())
			fn()
		}))
}
