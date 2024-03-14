package client

import (
	"reflect"

	"chat2-client-go/protos"
)

// 定义用到的类型

var (
	TypePtrNormalAck    = reflect.TypeOf(&protos.NormalAck{})
	TypePtrQueryGateAck = reflect.TypeOf(&protos.QueryGateAck{})
)
