package service

const (
	// SystemAPI
	// 所有service都要支持的
	SystemAPI = "__sys__"
)

//	一般会有一个actor的service提供实际服务
//	nodeService实际上负责，配置读取等
type INodeService interface {
	AddComponent(name string, comp IComponent)
	RemoveComponent(name string)
	GetComponent(name string) IComponent
}

//	嵌入nodeservice的必须实现此接口
type INodeServiceOwner interface {
	GetNodeService() *NodeService
}

// IComponent
// 可以加入到nodeservice中
type IComponent interface {
	Init(service INodeService)
	OnAdd()
	OnRemove()
}

// IServiceCreator
// 构建器
type IServiceCreator interface {
	Create(name string)
}

type IServiceFactory interface {
	Register(typeName string, creator IServiceCreator)
	Create(typeName string, name string)
}

// ICtrlCmdListener 控制消息处理者, 比如退休
type ICtrlCmdListener interface {
	Handler(cmd string) string
}
