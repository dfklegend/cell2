package module

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dfklegend/cell2/baseapp/interfaces"
	"github.com/dfklegend/cell2/utils/runservice"
)

type TestEnv struct {
	result string
}

func (t *TestEnv) Reset() {
	t.result = ""
}

var (
	theEnv = &TestEnv{}
)

// --------
type Module1 struct {
	*BaseModule
}

func NewModule1() *Module1 {
	return &Module1{
		NewBaseModule(),
	}
}

func (b *Module1) Start(next interfaces.FuncWithSucc) {
	log.Printf("Module1.Start\n")
	theEnv.result += "M1.Start"
	b.RunService.GetTimerMgr().After(time.Millisecond, func(args ...interface{}) {
		log.Printf("Module1.Start over\n")
		theEnv.result += "M1.StartOver"

		next(true)

	})
}

func (b *Module1) Stop(next interfaces.FuncWithSucc) {
	log.Printf("Module1.Stop\n")

	theEnv.result += "M1.Stop"
	b.RunService.GetTimerMgr().After(time.Millisecond, func(args ...interface{}) {
		log.Printf("Module1.Stop Over\n")
		theEnv.result += "M1.StopOver"

		next(true)
	})
}

type Module2 struct {
	*BaseModule
}

func NewModule2() *Module2 {
	return &Module2{
		NewBaseModule(),
	}
}

func (b *Module2) Start(next interfaces.FuncWithSucc) {
	log.Printf("Module2.Start\n")
	theEnv.result += "M2.Start"
	b.RunService.GetTimerMgr().After(time.Millisecond, func(args ...interface{}) {
		log.Printf("Module2.Start Over\n")
		theEnv.result += "M2.StartOver"

		next(true)
	})
}

func (b *Module2) Stop(next interfaces.FuncWithSucc) {
	log.Printf("Module2.Stop\n")
	theEnv.result += "M2.Stop"
	b.RunService.GetTimerMgr().After(time.Millisecond, func(args ...interface{}) {
		log.Printf("Module2.Stop over\n")
		theEnv.result += "M2.StopOver"
		next(true)
	})
}

type Module3 struct {
	*BaseModule
}

func NewModule3() *Module3 {
	return &Module3{
		NewBaseModule(),
	}
}

func (b *Module3) Start(next interfaces.FuncWithSucc) {
	log.Printf("Module3.Start\n")
	theEnv.result += "M3.Start"
	b.RunService.GetTimerMgr().After(time.Millisecond, func(args ...interface{}) {
		log.Printf("Module3.Start Over\n")
		theEnv.result += "M3.StartOver"

		next(false)
	})
}

func (b *Module3) Stop(next interfaces.FuncWithSucc) {
	log.Printf("Module3.Stop\n")
	theEnv.result += "M3.Stop"
	b.RunService.GetTimerMgr().After(time.Millisecond, func(args ...interface{}) {
		log.Printf("Module3.Stop over\n")
		theEnv.result += "M3.StopOver"
		next(true)
	})
}

// --------
//	Start顺序
//	Stop顺序
func Test_RunOrder(t *testing.T) {
	theEnv.Reset()
	service := runservice.NewStandardRunService("test")
	service.Start()

	ms := NewModList()

	m1 := NewModule1()
	m1.Init(service)
	ms.AddModule(m1)

	m2 := NewModule2()
	m2.Init(service)
	ms.AddModule(m2)

	ms.Start(func(succ bool) {
		log.Printf("start over:%v\n", succ)
	})

	time.Sleep(time.Second)

	log.Printf("%v\n", theEnv.result)
	assert.Equal(t, "M1.StartM1.StartOverM2.StartM2.StartOver", theEnv.result)

	theEnv.Reset()
	ms.Stop(func(succ bool) {
		log.Printf("stop over:%v\n", succ)
	})

	time.Sleep(time.Second)
	log.Printf("%v\n", theEnv.result)
	assert.Equal(t, "M2.StopM2.StopOverM1.StopM1.StopOver", theEnv.result)

	service.Stop()
}

// 测试中间失败
func Test_ErrBreak(t *testing.T) {
	theEnv.Reset()
	service := runservice.NewStandardRunService("test")
	service.Start()

	ms := NewModList()

	m1 := NewModule1()
	m1.Init(service)
	ms.AddModule(m1)

	m3 := NewModule3()
	m3.Init(service)
	ms.AddModule(m3)

	m2 := NewModule2()
	m2.Init(service)
	ms.AddModule(m2)

	ms.Start(func(succ bool) {
		log.Printf("start over:%v\n", succ)
	})

	time.Sleep(time.Second)

	log.Printf("%v\n", theEnv.result)
	assert.Equal(t, "M1.StartM1.StartOverM3.StartM3.StartOver", theEnv.result)

	theEnv.Reset()
	ms.Stop(func(succ bool) {
		log.Printf("stop over:%v\n", succ)
	})

	time.Sleep(time.Second)
	log.Printf("%v\n", theEnv.result)
	assert.Equal(t, "M2.StopM2.StopOverM3.StopM3.StopOverM1.StopM1.StopOver", theEnv.result)

	service.Stop()
}
