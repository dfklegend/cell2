package waterfall

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dfklegend/cell2/utils/common"
)

type TestEnv struct {
	resultErr     bool
	result        int
	lastRoutine   int
	inSameRoutine bool
}

func (te *TestEnv) Reset() {
	te.resultErr = false
	te.result = -1
	te.lastRoutine = -1
	te.inSameRoutine = true
}

func (te *TestEnv) UpdateRoutine() {
	curRoutineId := common.GetRoutineID()
	if te.lastRoutine == curRoutineId {
		return
	}

	if te.lastRoutine != -1 {
		te.inSameRoutine = false
	}
	te.lastRoutine = curRoutineId
}

var testEnv = &TestEnv{}

func Test_Simple(t *testing.T) {
	fmt.Println("---- Test_Simple ----")
	testEnv.Reset()
	go testDo(Simple)
	time.Sleep(time.Second * 2)

	assert.Equal(t, false, testEnv.resultErr)
	assert.Equal(t, 3, testEnv.result)
	assert.Equal(t, false, testEnv.inSameRoutine)
	fmt.Println("--------")
}

func Test_SimpleFail(t *testing.T) {
	fmt.Println("---- Test_SimpleFail ----")
	testEnv.Reset()
	go testFailDo(Simple)
	time.Sleep(time.Second * 2)

	assert.Equal(t, true, testEnv.resultErr)
	assert.Equal(t, 3, testEnv.result)
	assert.Equal(t, true, testEnv.inSameRoutine)
	fmt.Println("--------")
}

func Test_Go(t *testing.T) {
	fmt.Println("---- Test_Go ----")

	testEnv.Reset()
	go testDo(ExecAndWait)
	time.Sleep(time.Second * 2)

	assert.Equal(t, false, testEnv.resultErr)
	assert.Equal(t, 3, testEnv.result)
	assert.Equal(t, true, testEnv.inSameRoutine)
	fmt.Println("--------")
}

func Test_GoFail(t *testing.T) {
	fmt.Println("---- show waterfall go ----")

	testEnv.Reset()
	go testFailDo(ExecAndWait)
	time.Sleep(time.Second * 3)

	assert.Equal(t, true, testEnv.resultErr)
	assert.Equal(t, 3, testEnv.result)
	assert.Equal(t, true, testEnv.inSameRoutine)
	fmt.Println("--------")
}

// 一个固定流程
func testDo(fn func(tasks []Task, final FinalCallback)) {
	fmt.Println("begin routine:", common.GetRoutineID())
	fn([]Task{func(callback Callback, args ...interface{}) {
		//fmt.Println("task0 routine:", common.GetRoutineID())
		testEnv.UpdateRoutine()
		callback(false, 1, 2)
	}, func(callback Callback, args ...interface{}) {
		//fmt.Println("task1 routine:", common.GetRoutineID())
		x, _ := args[0].(int)
		y, _ := args[1].(int)
		testEnv.UpdateRoutine()
		callback(false, x+y)
	}, func(callback Callback, args ...interface{}) {
		fmt.Println("task2 routine:", common.GetRoutineID())
		go func() {
			fmt.Println("enter newgo:", common.GetRoutineID())
			time.Sleep(time.Second * 1)
			x, _ := args[0].(int)
			callback(false, x)
			fmt.Println("newgo over:", common.GetRoutineID())
		}()

	}}, func(err bool, args ...interface{}) {
		//fmt.Println("final routine:", common.GetRoutineID())
		fmt.Println(args...)
		testEnv.UpdateRoutine()

		testEnv.resultErr = err
		testEnv.result, _ = args[0].(int)
	})
}

//	中途出错，跳出
func testFailDo(fn func(tasks []Task, final FinalCallback)) {
	fmt.Println("begin routine:", common.GetRoutineID())
	fn([]Task{func(callback Callback, args ...interface{}) {
		//fmt.Println("task0 routine:", common.GetRoutineID())
		testEnv.UpdateRoutine()
		callback(false, 1, 2)
	}, func(callback Callback, args ...interface{}) {
		x, _ := args[0].(int)
		y, _ := args[1].(int)
		//fmt.Println("task1 routine:", common.GetRoutineID())
		testEnv.UpdateRoutine()
		callback(true, x+y)
	}, func(callback Callback, args ...interface{}) {
		fmt.Println("task2 routine:", common.GetRoutineID())
		go func() {
			fmt.Println("enter newgo:", common.GetRoutineID())
			time.Sleep(time.Second * 2)
			x, _ := args[0].(int)
			callback(false, x)
			fmt.Println("newgo over:", common.GetRoutineID())
		}()

	}}, func(err bool, args ...interface{}) {
		//fmt.Println("final routine:", common.GetRoutineID())
		fmt.Println(args...)
		testEnv.UpdateRoutine()

		testEnv.resultErr = err
		testEnv.result, _ = args[0].(int)
	})
}
