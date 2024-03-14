package runservice

import (
	"reflect"

	"github.com/dfklegend/cell2/utils/event"
	"github.com/dfklegend/cell2/utils/sche"
	"github.com/dfklegend/cell2/utils/timer"
)

// 	标准RunService，方便使用
// 	拥有timer
// 	事件中心
type StandardRunService struct {
	*RunService
	TimerMgr *timer.Mgr
	// 异步事件中心
	EventCenter *event.LocalEventCenter
}

func NewStandardRunService(name string) *StandardRunService {
	return &StandardRunService{
		RunService:  NewRunService(name),
		TimerMgr:    timer.NewTimerMgr(),
		EventCenter: event.NewLocalEventCenter(true),
	}
}

func (s *StandardRunService) GetTimerMgr() *timer.Mgr {
	return s.TimerMgr
}

func (s *StandardRunService) GetEventCenter() *event.LocalEventCenter {
	return s.EventCenter
}

func (s *StandardRunService) Start() {
	s.RunService.Start()
	// do something else
	s.addTimerSelector()
	s.addEventSelector()
}

func (s *StandardRunService) Stop() {
	// do something else
	s.TimerMgr.Stop()
	s.EventCenter.Clear()

	s.RunService.Stop()
}

func (s *StandardRunService) addTimerSelector() {
	mgr := s.TimerMgr
	selector := s.selector
	selector.AddSelector("timer", sche.NewFuncSelector(reflect.ValueOf(mgr.GetQueue()),
		func(v reflect.Value, recvOk bool) {
			if !recvOk {
				return
			}

			t := v.Interface().(*timer.Obj)
			mgr.Do(t)
		}))
}

func (s *StandardRunService) addEventSelector() {
	ec := s.EventCenter
	selector := s.selector
	selector.AddSelector("event", sche.NewFuncSelector(reflect.ValueOf(ec.GetChanEvent()),
		func(v reflect.Value, recvOk bool) {
			if !recvOk {
				return
			}

			e := v.Interface().(*event.EObj)
			ec.DoEvent(e)
		}))
}
