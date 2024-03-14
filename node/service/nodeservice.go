package service

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"

	as "github.com/dfklegend/cell2/actorex/service"
	messages "github.com/dfklegend/cell2/actorex/service/servicemsgs"
	"github.com/dfklegend/cell2/node/config"
	"github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/logger/interfaces"
)

// NodeService
// 节点服务
//   能定义客户端协议处理(非前端自动转发)
//   能提供服务器组内RPC
//   使用服务发现来进行服务维护
// need send StartServiceCmd after created
// 可以设置一些参数传递进来
type NodeService struct {
	*as.Service

	Name string
	info *config.ServiceInfo

	comps map[string]IComponent

	// 处理控制
	ctrlCmdListener ICtrlCmdListener
	// 设置后，通过nodeservice可以获取NodeService嵌入者
	owner INodeServiceOwner

	// log时，自动输出服务id
	log interfaces.Logger
}

func NewService() *NodeService {
	s := &NodeService{
		Service: as.NewService(),
		comps:   make(map[string]IComponent),
		log:     logger.L,
	}

	s.Service.InitReqReceiver(s)
	return s
}

func (n *NodeService) SetOwner(owner INodeServiceOwner) {
	n.owner = owner
}

func (n *NodeService) GetOwner() INodeServiceOwner {
	return n.owner
}

func (n *NodeService) SetCtrlCmdListener(listener ICtrlCmdListener) {
	n.ctrlCmdListener = listener
}

func (n *NodeService) GetCtrlCmdListener() ICtrlCmdListener {
	return n.ctrlCmdListener
}

func (n *NodeService) GetLogger() interfaces.Logger {
	return n.log
}

func (n *NodeService) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *StartServiceCmd:
		n.onStart(msg.Name, msg.Info)
		//log.Printf("set name:%v \n", n.Name)
	}
	n.Service.Receive(ctx)
}

func (n *NodeService) onStart(name string, info *config.ServiceInfo) {
	n.Name = name
	n.info = info

	n.log = logger.L.WithField("service", name)
	n.log.Infof("service start")
}

func (n *NodeService) GetServiceType() string {
	return n.info.Type
}

func (n *NodeService) ReceiveRequest(ctx actor.Context, request *messages.ServiceRequest, rawMsg interface{}) {
}

func (n *NodeService) AddComponent(name string, comp IComponent) {
	if n.comps[name] != nil {
		log.Printf("already has component with name: %v\n", name)
		return
	}
	comp.Init(n)
	n.comps[name] = comp

	comp.OnAdd()
}

func (n *NodeService) GetComponent(name string) IComponent {
	return n.comps[name]
}

func (n *NodeService) RemoveComponent(name string) {
	comp := n.comps[name]
	if comp == nil {
		return
	}
	comp.OnRemove()
	delete(n.comps, name)
}
