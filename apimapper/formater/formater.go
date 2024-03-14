package formater

import (
	"reflect"

	api "github.com/dfklegend/cell2/apimapper"
)

var defaultFormater api.IAPIFormatter = &DefaultFormater{}

func GetDefaultFormater() api.IAPIFormatter {
	return defaultFormater
}

// 	DefaultFormater 缺省的rpc格式
type DefaultFormater struct {
}

func (d *DefaultFormater) IsValidMethod(method reflect.Method) bool {
	return d.isValidRequest(method) || d.isValidNotify(method)
}

// 	handler函数 4个参数
// 		receiver, context, pointer of msg, cb
func (d *DefaultFormater) isValidRequest(method reflect.Method) bool {
	mt := method.Type
	// Method must be exported.
	if method.PkgPath != "" {
		return false
	}

	// Method needs three ins: receiver, context, pointer, cb.
	if mt.NumIn() != 4 {
		//log.Printf("%v must had %v args\n", method.Name, 3)
		return false
	}

	// Method needs one outs: error
	//if mt.NumOut() != 1 || mt.Out(0) != api.TypeOfError {
	//	return false
	//}

	// 参数1是IContext
	if t1 := mt.In(1); t1.Kind() != reflect.Ptr || !t1.Implements(api.TypeOfContext) {
		return false
	}

	// 参数2 消息体
	if mt.In(2).Kind() != reflect.Ptr {
		return false
	}

	if mt.In(3).Kind() != reflect.Func {
		return false
	}
	return true
}

// 	handler函数 3个参数
// 		receiver, context, pointer of msg
func (d *DefaultFormater) isValidNotify(method reflect.Method) bool {
	mt := method.Type
	// Method must be exported.
	if method.PkgPath != "" {
		return false
	}

	// Method needs three ins: receiver, context, pointer, cb.
	if mt.NumIn() != 3 {
		//log.Printf("%v must had %v args\n", method.Name, 3)
		return false
	}

	// Method needs one outs: error
	//if mt.NumOut() != 1 || mt.Out(0) != api.TypeOfError {
	//	return false
	//}

	// 参数1是IContext
	if t1 := mt.In(1); t1.Kind() != reflect.Ptr || !t1.Implements(api.TypeOfContext) {
		return false
	}

	// 参数2 消息体
	if mt.In(2).Kind() != reflect.Ptr {
		return false
	}

	return true
}
