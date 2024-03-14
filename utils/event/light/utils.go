package light

import (
	"reflect"
)

// Bind接口避免书写错误

func BindEventWithReceiver(bind bool, center *EventCenter, eventName string, receiver any, cb CBFunc, args ...interface{}) {
	if bind {
		center.SubscribeWithReceiver(eventName, receiver, cb, args...)
	} else {
		center.UnsubscribeWithReceiver(eventName, receiver, cb)
	}
}

func BindEvent(bind bool, center *EventCenter, eventName string, cb CBFunc, args ...interface{}) {
	if bind {
		center.Subscribe(eventName, cb, args...)
	} else {
		center.Unsubscribe(eventName, cb)
	}
}

func toAny(in any) any {
	return in
}

func toPointer(in any) uintptr {
	if in == nil {
		return 0
	}
	return reflect.ValueOf(in).Pointer()
}

func compareFunc(a any, b any) bool {
	aPointer := toPointer(a)
	bPointer := toPointer(b)
	return aPointer == bPointer
}
