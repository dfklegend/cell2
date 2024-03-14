package waterfall

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dfklegend/cell2/utils/common"
	scheP "github.com/dfklegend/cell2/utils/sche"
)

func Test_Sche_Simple(t *testing.T) {
	sche := scheP.NewSche()

	result := 0
	Sche(sche, []Task{
		func(callback Callback, args ...interface{}) {
			callback(false, 1)
		},
		func(callback Callback, args ...interface{}) {
			arg1 := args[0].(int)
			callback(false, arg1+2)
		},
	}, func(err bool, args ...interface{}) {
		result = args[0].(int)
	})

	go sche.Handler()
	time.Sleep(time.Second * 1)
	sche.Stop()

	assert.Equal(t, 3, result)
}

func Test_Sche_SimpleBuilder(t *testing.T) {
	sche := scheP.NewSche()

	result := 0
	NewBuilder(sche).Next(func(callback Callback, args ...interface{}) {
		callback(false, 1)
	}).Next(func(callback Callback, args ...interface{}) {
		arg1 := args[0].(int)
		callback(false, arg1+2)
	}).Final(func(err bool, args ...interface{}) {
		result = args[0].(int)
	}).Do()

	go sche.Handler()
	time.Sleep(time.Second * 1)
	sche.Stop()

	assert.Equal(t, 3, result)
}

func Test_Sche_Panic(t *testing.T) {
	sche := scheP.NewSche()

	result := 0
	Sche(sche, []Task{
		func(callback Callback, args ...interface{}) {
			callback(false, "hello")
		},
		func(callback Callback, args ...interface{}) {
			// will panic
			arg1 := args[0].(int)
			callback(false, arg1+2)
		},
	}, func(err bool, args ...interface{}) {
		result = -1
	})

	go sche.Handler()
	time.Sleep(time.Second * 1)
	sche.Stop()

	assert.Equal(t, 0, result)
}

func Test_ExecInOneRoutine(t *testing.T) {
	fmt.Println("---- Test_InOneRoutine ----")

	testEnv.Reset()
	sche := scheP.NewSche()
	fn := func(tasks []Task, final FinalCallback) {
		Sche(sche, tasks, final)
	}
	go testDo(fn)

	go func() {
		fmt.Println("taskInRoutine:", common.GetRoutineID())
		sche.Handler()
	}()
	time.Sleep(time.Second * 3)
	sche.Stop()

	assert.Equal(t, true, testEnv.inSameRoutine)
	fmt.Println("all over")
}
