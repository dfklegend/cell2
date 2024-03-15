package systems

import (
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/utils/event/light"

	mymsg "mmo/messages"
)

type ILogicPlayer interface {
	GetUId() int64

	SetDirt()
	GetNodeService() *service.NodeService
	GetEvents() *light.EventCenter

	PushMsg(route string, msg interface{})

	GetSystem(name string) ISystem
}

// ISystem player身上的逻辑组织成子系统
// 代码拆分

// 创建对象时构建
// system
// system  load/save数据
// system 原则上不和其他system发生关联
// 可以通过事件协作
// system和player可以通过接口来操作

type ISystem interface {
	Init(player ILogicPlayer)

	OnCreate()
	OnDestroy()

	// PushInfoToClient 向客户端Push一些初始化信息
	PushInfoToClient()
	// OnEnterWorld 可以基于逻辑做一些数据初始化
	OnEnterWorld(switchLine bool)

	// InitData 创角时，数据初始化
	InitData()
	LoadData(info *mymsg.PlayerInfo)
	SaveData(info *mymsg.PlayerInfo)

	// Request 可以向具体system请求request
	// 定义一些子协议
	Request(cmd string, args []byte, cb func(ret []byte, errCode int32))
	HasCmd(cmd string) bool
}
