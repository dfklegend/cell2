package sche

import (
	"fmt"
	"testing"
	"time"

	"github.com/dfklegend/cell2/utils/common"

	"github.com/stretchr/testify/assert"
)

type Null struct{}

// 测试处理中产生了panic
func Test_Panic(t *testing.T) {
	fmt.Println("---- Test_Panic ----")

	sche := NewSche()

	var rtSche, rtCall int
	go func() {
		fmt.Printf("%v call routine:%d\n", time.Now(), common.GetRoutineID())
		rtCall = common.GetRoutineID()
		sche.Post(func() {
			fmt.Printf("exec in routine:%d\n", common.GetRoutineID())
			rtSche = common.GetRoutineID()
			//var x interface{} = 7
			//return x.(struct{}), nil
			//var x = Null{}
			time.Sleep(1 * time.Second)
			panic("xxxx")
			//return 7, nil
		})

		fmt.Printf("%v call routine over:%d\n", time.Now(), common.GetRoutineID())
	}()

	go sche.Handler()

	time.Sleep(2 * time.Second)

	// 检查结果
	assert.NotEqual(t, rtSche, rtCall, "应该不在同一个routine")
}

func Test_ScheClose(t *testing.T) {
	fmt.Println("---- Test_ScheClose ----")

	sche := NewSche()
	fmt.Printf("%v main routine:%d\n", time.Now(), common.GetRoutineID())
	_ = sche.Post(func() {
		fmt.Printf("%v call routine over:%d\n", time.Now(), common.GetRoutineID())
	})
	go sche.Handler()

	time.Sleep(1 * time.Second)

	sche.Stop()
	_ = sche.Post(func() {
		fmt.Printf("%v call routine over:%d\n", time.Now(), common.GetRoutineID())
	})

	time.Sleep(3 * time.Second)
}

func Test_ScheFullBlock(t *testing.T) {
	// 不会阻塞处理事件携程
	selfBlockDefend = true
	assert.Equal(t, QueueSize+10, doTestScheFullBlock())

	// 关闭后，post阻塞了，一次执行不了
	selfBlockDefend = false
	assert.Equal(t, 0, doTestScheFullBlock())
	selfBlockDefend = false
}

func doTestScheFullBlock() int {
	sche := NewSche()

	running := true
	num := 0
	go func() {
		for i := 0; i < QueueSize+10; i++ {
			sche.Post(func() {})
		}

		for running {
			select {
			case <-sche.GetChanTask():
				num++
			case <-time.After(time.Second):
				running = false
				break
			}
		}
		fmt.Printf("got %v \n", num)
	}()

	wait := 5
	for running && wait > 0 {
		time.Sleep(time.Second)
		wait--
	}

	running = false

	sche.Stop()
	time.Sleep(time.Second)
	return num
}
