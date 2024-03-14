package nodectrl

import (
	"fmt"
	"strconv"

	"github.com/dfklegend/cell2/nodectrl/define"
	"github.com/dfklegend/cell2/utils/jsonutils"
	l "github.com/dfklegend/cell2/utils/logger"
)

// ----
func initCmds(cmds *Cmds) {
	cmds.Register("stat", &StatCmd{})
	cmds.Register("retire", &RetireCmd{})
	cmds.Register("exit", &ExitCmd{})
	//web master
	cmds.Register("web_nodes", &WebCmdNodes{})
	cmds.Register("web_retire", &WebCmdRetire{})
	cmds.Register("web_exit", &WebCmdExit{})
}

type NodeStatus struct {
	Node     string            `json:"node"`
	Service  string            `json:"service"`
	Status   string            `json:"status"`
	Time     string            `json:"time"`
	NodeData map[string]string `json:"nodeData"`
}

// ----

type StatCmd struct {
}

func (c *StatCmd) Handle(n *NodeCtrl, args string) string {
	stat := n.GetStatStr()
	l.L.Infof("%v", stat)
	return stat
}

// ----

type RetireCmd struct {
}

func (c *RetireCmd) Handle(n *NodeCtrl, args string) string {
	if !n.IsState(define.Working) && !n.IsState(define.Retiring) {
		l.L.Errorf("beginRetire failed, error state: %v", n.state)
		return fmt.Sprintf("beginRetire failed, error state: %v", define.GetStateName(n.state))
	}

	if !n.retireSupport {
		l.L.Errorf("some service not support retire")
		n.dumpAllRetireUnSupport()
		return "some service not support retire"
	}

	// 开始退休
	n.setState(define.Retiring)
	l.L.Infof("node -> Retiring")
	// 要求各service开始退休
	n.SendCmdToAllSelfServices("retire")
	return "ok"
}

// ----

type ExitCmd struct {
}

func (c *ExitCmd) Handle(n *NodeCtrl, args string) string {
	if !n.IsState(define.Retired) {
		l.L.Errorf("beginExit failed, error state: %v", n.state)
		return fmt.Sprintf("beginExit failed, error state: %v", define.GetStateName(n.state))
	}

	n.setState(define.Exiting)
	n.node.StopNode(func(succ bool) {
		l.L.Infof("NodeCtrl.beginExit ret: %v", succ)
		if succ {
			n.setState(define.Exited)
		}
	})
	return "ok"
}

// WebCmdNodes 节点状态
type WebCmdNodes struct {
}

//TODO 测试数据
var online = 0

func (w *WebCmdNodes) Handle(n *NodeCtrl, args string) string {

	online++
	n.SetNodeData("online", strconv.Itoa(online))

	return jsonutils.Marshal(&NodeStatus{
		Node:     n.GetNodeId(),
		Status:   define.GetStateName(n.state),
		Time:     fmt.Sprintf("%v", n.GetPassed()/1000),
		Service:  jsonutils.Marshal(n.services),
		NodeData: map[string]string{"online": n.GetNodeData("online")},
	})
}

// WebCmdRetire 节点退休
type WebCmdRetire struct {
}

func (w *WebCmdRetire) Handle(n *NodeCtrl, args string) string {
	if !n.IsState(define.Working) && !n.IsState(define.Retiring) {
		l.L.Errorf("beginRetire failed, error state: %v", n.state)
		return fmt.Sprintf("beginRetire failed, error state: %v", define.GetStateName(n.state))
	}

	if !n.retireSupport {
		l.L.Errorf("some service not support retire")
		n.dumpAllRetireUnSupport()
		return "some service not support retire"
	}

	// 开始退休
	n.setState(define.Retiring)
	l.L.Infof("node -> Retiring")
	// 要求各service开始退休
	n.SendCmdToAllSelfServices("retire")
	return "ok"
}

// WebCmdExit 节点退出
type WebCmdExit struct {
}

func (w *WebCmdExit) Handle(n *NodeCtrl, args string) string {
	if !n.IsState(define.Retired) {
		l.L.Errorf("beginExit failed, error state: %v", n.state)
		return fmt.Sprintf("beginExit failed, error state: %v", define.GetStateName(n.state))
	}

	n.setState(define.Exiting)
	n.node.StopNode(func(succ bool) {
		l.L.Infof("NodeCtrl.beginExit ret: %v", succ)
		if succ {
			n.setState(define.Exited)
		}
	})
	return "ok"
}
