package module

import (
	"time"

	"github.com/dfklegend/cell2/utils/runservice"
)

type BaseModule struct {
	RunService *runservice.StandardRunService
}

func NewBaseModule() *BaseModule {
	return &BaseModule{}
}

func (b *BaseModule) Init(rs *runservice.StandardRunService) {
	b.RunService = rs
}

func (b *BaseModule) Start(next func()) {
	b.RunService.GetTimerMgr().After(time.Millisecond, func(args ...interface{}) {
		next()
	})
}

func (b *BaseModule) Stop(next func()) {
	b.RunService.GetTimerMgr().After(time.Millisecond, func(args ...interface{}) {
		next()
	})
}
