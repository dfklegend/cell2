package apientry

import (
	"errors"
	"fmt"
	"reflect"

	api "github.com/dfklegend/cell2/apimapper"
	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/serialize"
)

// CallWithSerialize
// 序列化入参
// 必返回cb
func CallWithSerialize(c *APICollection, ctx api.IContext, route string, data []byte,
	cbFunc HandlerCBFunc, serializer serialize.Serializer) {
	if serializer == nil {
		CheckInvokeCBFunc(cbFunc, errors.New(fmt.Sprintf("serializer is nil: %v", route)), nil)
		return
	}
	argType := c.GetArgType(route)
	if argType == nil {
		CheckInvokeCBFunc(cbFunc, errors.New(fmt.Sprintf("no method: %v", route)), nil)
		return
	}

	arg := reflect.New(argType.Elem()).Interface()
	err := serializer.Unmarshal(data, arg)
	if err != nil {
		l.Log.Infof("arg Unmarshal [%v] failed: %v\n", string(data), err)
		CheckInvokeCBFunc(cbFunc, err, nil)
		return
	}

	c.Call(ctx, route, arg, cbFunc)
}
