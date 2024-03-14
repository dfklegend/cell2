package waterfall

import (
	"fmt"
	"log"

	"github.com/dfklegend/cell2/utils/common"
)

type Callback func(err bool, args ...interface{})
type FinalCallback func(err bool, args ...interface{})
type Task func(callback Callback, args ...interface{})

// 实现一个类似nodejs waterfall

//	简版
//	测试函数指针等
//	just callback directly
func Simple(tasks []Task, final FinalCallback) {
	// 执行索引
	cursor := 0
	size := len(tasks)

	var callback Callback = nil

	exec := func(index int, args ...interface{}) {
		tasks[index](callback, args...)
	}

	gonext := func(args ...interface{}) {
		cursor++
		if cursor < size {
			exec(cursor, args...)
		} else {
			//args = append([]interface{}{false}, args...)
			final(false, args...)
		}
	}

	callback = func(err bool, args ...interface{}) {
		//fmt.Printf("callback:%v\n", util.GetRoutineID())
		if err {
			// go final
			final(true, args...)
			return
		}

		// go next
		gonext(args...)
	}

	exec(0)
}

// 	每个Task都在routine中执行
//	阻塞等待执行完毕
// 	every task will execute in one routine(caller routine)
func ExecAndWait(tasks []Task, final FinalCallback) {
	// 执行索引
	cursor := 0
	size := len(tasks)
	//
	taskNext := 0
	taskFinal := 1

	chanNext := make(chan int, 1)

	var callback Callback = nil
	var gofinal Callback = nil
	var curArgs []interface{}
	var curErr = false

	exec := func(index int, args ...interface{}) {
		tasks[index](callback, args...)
	}

	gonext := func(args ...interface{}) {
		curArgs = args
		chanNext <- taskNext
	}

	donext := func(args ...interface{}) {
		cursor++
		if cursor < size {
			exec(cursor, args...)
		} else {
			// error: false
			gofinal(false, args...)
		}
	}

	gofinal = func(err bool, args ...interface{}) {
		curArgs = args
		curErr = err
		chanNext <- taskFinal
	}

	dofinal := func(err bool, args ...interface{}) {
		final(err, args...)
		close(chanNext)
	}

	// 不在主routine中
	//
	callback = func(err bool, args ...interface{}) {
		log.Printf("callback:%v  %v\n", common.GetRoutineID(), args)
		if err {
			// go final
			//args = append([]interface{}{true}, args...)

			log.Printf("%v\n", args)
			gofinal(true, args...)
			return
		}

		// go next
		gonext(args...)
	}

	// 启动
	exec(0)

	// 问题，会阻塞调用者的routine
	// wait for next task
	for {
		// 0: next
		// 1: final
		data, ok := <-chanNext
		if !ok {
			// over
			fmt.Println("waterfall over")
			return
		}

		if data == taskNext {
			donext(curArgs...)
		} else {
			dofinal(curErr, curArgs...)
		}
	}
}
