package waterfall

import (
	scheP "github.com/dfklegend/cell2/utils/sche"
)

// 实现一个类似nodejs waterfall
// 每个回调都在调用routine中调用

type Chain struct {
	// 任务列表
	tasks []Task
	// 最后的任务
	final  FinalCallback
	cursor int
	// 执行
	callbackFunc Callback
}

func (c *Chain) invokeTask(index int, args ...interface{}) {
	c.tasks[index](c.callbackFunc, args...)
}

func (c *Chain) tryExec(index int, args ...interface{}) {
	if index < len(c.tasks) {
		c.invokeTask(index, args...)
	} else {
		c.invokeFinal(false, args...)
	}
}

func (c *Chain) next(args ...interface{}) {
	c.cursor++
	c.tryExec(c.cursor, args...)
}

func (c *Chain) invokeFinal(err bool, args ...interface{}) {
	c.final(err, args...)
}

func (c *Chain) invokeCallback(err bool, args ...interface{}) {
	if err {
		c.invokeFinal(true, args...)
		return
	}

	// next
	c.next(args...)
}

// 	一般远程任务，传入的回调函数，可能在其他的执行环境内触发
// 	保证Task和Final都在sche内驱动执行
func Sche(sche *scheP.Sche, tasks []Task, final FinalCallback) *Chain {
	thisChain := &Chain{
		tasks: tasks,
		final: final,
	}

	thisChain.callbackFunc = func(err bool, args ...interface{}) {
		// 去到环境中执行
		sche.Post(func() {
			thisChain.invokeCallback(err, args...)
		})
	}

	sche.Post(func() {
		thisChain.tryExec(0)
	})
	return thisChain
}

//	Builder of sche
//	NewBuilder().Next().Next().Final().Do()
type Builder struct {
	sche  *scheP.Sche
	tasks []Task
	final FinalCallback
}

func NewBuilder(s *scheP.Sche) *Builder {
	return &Builder{
		sche:  s,
		tasks: make([]Task, 0),
	}
}

func (b *Builder) Next(task Task) *Builder {
	b.tasks = append(b.tasks, task)
	return b
}

func (b *Builder) Final(final FinalCallback) *Builder {
	b.final = final
	return b
}

func (b *Builder) Do() {
	Sche(b.sche, b.tasks, b.final)
}
