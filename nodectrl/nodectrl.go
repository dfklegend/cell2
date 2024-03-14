package nodectrl

import (
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/actorex/service"
	appdef "github.com/dfklegend/cell2/node/app/define"
	"github.com/dfklegend/cell2/node/builtin/msgs"
	"github.com/dfklegend/cell2/node/config"
	"github.com/dfklegend/cell2/nodectrl/define"
	"github.com/dfklegend/cell2/utils/common"
	l "github.com/dfklegend/cell2/utils/logger"
)

type ServiceState struct {
	State         define.NodeState
	RetireSupport bool
}

// NodeCtrl
// 节点控制器
// 负责节点退休，退出
type NodeCtrl struct {
	NodeId  string
	node    appdef.INodeApp
	service *AdminService
	state   define.NodeState
	cmds    *Cmds
	// 启动时间
	timeStart int64
	adminPID  *actor.PID

	services      map[string]*ServiceState
	retireSupport bool

	nodeData map[string]string
}

func NewNodeCtrl() *NodeCtrl {
	return &NodeCtrl{
		cmds:          NewCmds(),
		state:         define.Working,
		timeStart:     common.NowMs(),
		services:      make(map[string]*ServiceState),
		retireSupport: false,
		nodeData:      map[string]string{},
	}
}

func (n *NodeCtrl) Start(app appdef.INodeApp) {
	n.node = app
	n.adminPID = CreateAdminService(app.GetActorSystem(), func(s service.IService) {
		admin, _ := s.(*AdminService)
		n.service = admin
		admin.SetCtrl(n)
	})
	l.L.Infof("node ctrl started")
	n.makeServices()
	n.delayCheckRetireSupport()
	initCmds(n.cmds)
}

func (n *NodeCtrl) setState(state define.NodeState) {
	n.state = state
	n.node.UpdateNodeState(int(n.state))
}

func (n *NodeCtrl) GetState() define.NodeState {
	return n.state
}

func (n *NodeCtrl) IsState(state define.NodeState) bool {
	return n.state == state
}

func (n *NodeCtrl) GetPassed() int64 {
	return common.NowMs() - n.timeStart
}

func (n *NodeCtrl) GetAdmin() *actor.PID {
	return n.adminPID
}

func (n *NodeCtrl) makeServices() {
	n.node.FilterSelfServices(func(name string, cfg *config.ServiceInfo) {
		n.services[name] = &ServiceState{
			State:         define.Working,
			RetireSupport: false,
		}
	})
}

func (n *NodeCtrl) delayCheckRetireSupport() {
	// 探测一下是否支持"retire"策略
	n.service.GetRunService().GetTimerMgr().After(3*time.Second, func(args ...interface{}) {
		n.checkRetireSupport()
	})
}

func (n *NodeCtrl) checkRetireSupport() {
	for k, v := range n.services {
		n.queryRetire(k, v)
	}
}

func (n *NodeCtrl) queryRetire(name string, v *ServiceState) {
	n.sendCmd(name, "queryretire", func(result string) {
		if "ok" == result {
			v.RetireSupport = true

			l.L.Infof("service :%v support retire", name)
			n.retireSupport = n.checkAllRetireSupport()
		}
	})
}

func (n *NodeCtrl) checkAllRetireSupport() bool {
	for _, v := range n.services {
		if !v.RetireSupport {
			return false
		}
	}
	return true
}

func (n *NodeCtrl) dumpAllRetireUnSupport() {
	for k, v := range n.services {
		if !v.RetireSupport {
			l.L.Infof("%v not support retire", k)
		}
	}
}

func (n *NodeCtrl) GetStatStr() string {
	return fmt.Sprintf("state: %v, running passed: %v s", define.GetStateName(n.state), n.GetPassed()/1000)
}

func (n *NodeCtrl) ProcessCmd(cmd string) string {
	return n.cmds.Handle(n, cmd, "")
}

func (n *NodeCtrl) SendCmdToAllSelfServices(cmd string) {
	n.node.FilterSelfServices(func(name string, cfg *config.ServiceInfo) {
		n.sendCmd(name, cmd, nil)
	})
}

func (n *NodeCtrl) sendCmd(name string, cmd string, cb func(result string)) {
	pid := n.node.GetService(name)
	if pid == nil {
		l.L.Errorf("nodectrl can not find service: %v", name)
		return
	}

	n.service.RequestEx(pid, "ctrl.cmd", &msgs.CtrlCmd{
		Cmd: cmd,
	}, func(err error, raw interface{}) {
		if err != nil {
			l.L.Errorf("%v ctrl.cmd err: %v", name, err)
			return
		}
		if cb == nil {
			return
		}
		msg, _ := raw.(*msgs.CtrlCmdAck)
		cb(msg.Result)
	})
}

func (n *NodeCtrl) ProcessServiceCmd(name, cmd string) string {
	if cmd == "retired" {
		n.onServiceRetired(name)
		return "ok"
	}
	return "ok"
}

func (n *NodeCtrl) onServiceRetired(name string) {
	service := n.services[name]
	if service == nil {
		return
	}
	service.State = define.Retired

	if n.isAllServiceRetired() {
		// can exit
		n.setState(define.Retired)
		l.L.Infof("node -> retired")
	}
}

func (n *NodeCtrl) isAllServiceRetired() bool {
	for _, v := range n.services {
		if v.State != define.Retired {
			return false
		}
	}
	return true
}

func (n *NodeCtrl) GetNodeId() string {
	return n.NodeId
}

func (n *NodeCtrl) GetNodeData(key string) string {
	return n.nodeData[key]
}

func (n *NodeCtrl) SetNodeData(key string, value string) {
	n.nodeData[key] = value
}
