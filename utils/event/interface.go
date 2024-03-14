package event

/**
 * event包
 * 提供两级的事件中心
 * Global事件，可以推送到其他local
 * Local可以注册局部以及全局事件
 * Local事件，需使用者保证线程安全
 * (也可以使用SetLocalUseChan来确保局部事件也使用channel来
 * 拉取事件(能定制运行coroutine))
 * 并且需要自主拉取global事件 select ChanEvent
 *
 * 推荐事件参数还是定义成固定结构，便于维护
 *
 *
 * 非线程安全的更简洁的事件中心 light
 *
 */

type ChanEvent chan *EObj
type CBFunc func(args ...interface{})

type EObj struct {
	EventName string
	Args      []interface{}
}

type EListener struct {
	Id uint64
	// 自带参数
	Args []interface{}
	CB   CBFunc
}

// IGlobalEventCenter 全局的事件中心
type IGlobalEventCenter interface {
	Subscribe(eventName string, child ILocalEventCenter)
	Unsubscribe(eventName string, child ILocalEventCenter)
	Publish(eventName string, args ...interface{})
}

// ILocalEventCenter 本地事件中心，为本地服务，用chan拉取到保证cb的执行体
type ILocalEventCenter interface {
	// GetId 自身id
	GetId() uint64
	Clear()
	// Subscribe 本地事件
	Subscribe(eventName string, cb CBFunc, args ...interface{}) uint64
	Unsubscribe(eventName string, id uint64)
	// GSubscribe 全局事件
	GSubscribe(eventName string, cb CBFunc, args ...interface{}) uint64
	GUnsubscribe(eventName string, id uint64)

	GetChanEvent() ChanEvent
	DoEvent(e *EObj)

	Publish(eventName string, args ...interface{})
}
