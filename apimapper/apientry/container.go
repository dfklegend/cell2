package apientry

// 参考了pitaya

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	api "github.com/dfklegend/cell2/apimapper"
	"github.com/dfklegend/cell2/utils/common"
	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/serialize"
	"github.com/dfklegend/cell2/utils/serialize/json"
)

type HandlerCBFunc func(error, interface{})

//	Handler 对应一个函数
type Handler struct {
	// 对应的函数
	Method reflect.Method
	// 参数1 receiver
	ContextType reflect.Type
	// 参数3类型
	ArgType reflect.Type
	// 是否是request, 有cb
	IsRequest bool
}

// 	APIContainer 基于入口对象分析构建而成
//	缺省使用json serializer
type APIContainer struct {
	Name     string        // groupName of service
	Type     reflect.Type  // type of the receiver
	Receiver reflect.Value // receiver of methods for the service
	// "someMethod": Handler
	Handlers map[string]*Handler // registered methods
	Options  options             // options

	serializer   serialize.Serializer
	serializeRet bool
}

func NewContainer(entry api.IAPIEntry, opts ...Option) *APIContainer {
	s := &APIContainer{
		Type:         reflect.TypeOf(entry),
		Receiver:     reflect.ValueOf(entry),
		serializeRet: true,
	}

	// apply options
	for _, opt := range opts {
		opt(&s.Options)
	}

	if name := s.Options.groupName; name != "" {
		s.Name = name
	} else {
		s.Name = reflect.Indirect(s.Receiver).Type().Name()
		if s.Options.nameFunc != nil {
			s.Name = s.Options.nameFunc(s.Name)
		}
	}

	if s.Options.serializer != nil {
		s.serializer = s.Options.serializer
	} else {
		s.serializer = json.GetDefaultSerializer()
	}

	s.serializeRet = s.Options.serializeRet
	return s
}

// suitableMethods returns suitable methods of typ
func (c *APIContainer) suitableHandlerMethods(formatter api.IAPIFormatter, typ reflect.Type) map[string]*Handler {
	methods := make(map[string]*Handler)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mt := method.Type
		mn := method.Name

		// rewrite handler groupName
		if c.Options.nameFunc != nil {
			mn = c.Options.nameFunc(mn)
		}

		if formatter != nil && formatter.IsValidMethod(method) {
			methods[mn] = &Handler{
				Method:      method,
				ContextType: mt.In(1),
				ArgType:     mt.In(2),
				IsRequest:   mt.NumIn() == 4,
			}
		}
	}
	return methods
}

// ExtractHandler extract the set of methods from the
// receiver value which satisfy the following conditions:
// - exported method of exported type
func (c *APIContainer) ExtractHandler(formater api.IAPIFormatter) error {
	typeName := reflect.Indirect(c.Receiver).Type().Name()
	if typeName == "" {
		return errors.New("no service groupName for type " + c.Type.String())
	}
	if !isExported(typeName) {
		return errors.New("type " + typeName + " is not exported")
	}

	// Install the methods
	c.Handlers = c.suitableHandlerMethods(formater, c.Type)

	if len(c.Handlers) == 0 {
		str := ""
		// To help the user, see if a pointer receiver would work.
		method := c.suitableHandlerMethods(formater, reflect.PtrTo(c.Type))
		if len(method) != 0 {
			str = "type " + c.Name + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
		} else {
			str = "type " + c.Name + " has no exported methods of suitable type"
		}
		return errors.New(str)
	}

	for i := range c.Handlers {
		one := c.Handlers[i]
		l.Log.Infof("found Handler:%v %v", i, one.Method)
	}

	return nil
}

func (c *APIContainer) HasMethod(method string) bool {
	return c.Handlers[method] != nil
}

func (c *APIContainer) GetArgType(method string) reflect.Type {
	h := c.Handlers[method]
	if h == nil {
		return nil
	}
	return h.ArgType
}

// CallMethod
// 查找Handler
// 调用具体函数，无需序列化
/**
 * @param formater{IAPIFormatter} 用于提供接口类型差异化比如handler和remote
 * @param method{string} 方法名
 * @param msg{[]byte} 此接口负责反序列化成接口参数结构
 * @param cbFunc{HandlerCBFunc}
 *        cb(e error, outArg interface{})
 *        	e 错误
 *        	outArg 接口返回的对象 *
 * @param ext{interface{}} 见collection说明
 */
func (c *APIContainer) CallMethod(ctx api.IContext,
	method string,
	arg any, cbFunc HandlerCBFunc) {
	//log.Printf("enter: %v", method)
	handler := c.Handlers[method]
	if handler == nil {
		l.Log.Warnf("can not find method:%v", method)
		CheckInvokeCBFunc(cbFunc, fmt.Errorf("can not find method:%v", method), nil)
		return
	}

	// 参数列表
	// context
	var args []reflect.Value
	if handler.IsRequest {
		// 调用Request不带cb，也是合法
		// request
		args = []reflect.Value{c.Receiver,
			makeValueMaybeNil(handler.ContextType, ctx),
			//reflect.ValueOf(arg),
			makeValueMaybeNil(handler.ArgType, arg),
			reflect.ValueOf(cbFunc)}
	} else {
		// 如果调用notify,带一个cb,警告
		if cbFunc != nil {
			l.Log.Errorf("call notify with cb: %v.%v", c.Name, method)
			return
		}
		// notify
		args = []reflect.Value{c.Receiver,
			makeValueMaybeNil(handler.ContextType, ctx),
			//reflect.ValueOf(arg),
			makeValueMaybeNil(handler.ArgType, arg)}
	}

	SafeCall(handler, args, cbFunc)
	return
}

// SafeCall 捕获异常
func SafeCall(handler *Handler, args []reflect.Value, cbFunc HandlerCBFunc) {
	defer func() {
		if err := recover(); err != nil {
			l.E.Errorf("panic in handler.call:%v", err)

			stack := common.GetStackStr()
			log.Printf(stack)
			l.E.Infof(stack)

			CheckInvokeCBFunc(cbFunc, errors.New("panic in rpc"), nil)
		}
	}()

	handler.Method.Func.Call(args)
}

func (c *APIContainer) Dump() {
	L := l.L
	for k, _ := range c.Handlers {
		L.Infof("      func %v", k)
	}
}
