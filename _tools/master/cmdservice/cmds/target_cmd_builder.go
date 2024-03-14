package cmds

import (
	"fmt"
	"strings"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/actorex/service"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/cluster"
	"github.com/dfklegend/cell2/nodectrl/define"
)

// CmdTargetBuilder 对目标member执行cmd
type CmdTargetBuilder struct {
	service  *service.Service
	doFunc   func(b ICmdBuilder, pid *actor.PID, id string, member *cluster.Member)
	cbFin    func(string)
	makeFunc func([]string) string

	results []string
	waitNum int

	// 缺省命令
	cmd    string
	target string
}

func NewCmdTargetBuilder() *CmdTargetBuilder {
	return &CmdTargetBuilder{}
}

func (b *CmdTargetBuilder) WithCmd(cmd string, target string) *CmdTargetBuilder {
	b.cmd = cmd
	b.target = target
	b.doFunc = DefaultCmdDo
	b.makeFunc = DefaultTargetMakeResult
	return b
}

func (b *CmdTargetBuilder) WithMakeResult(f func(results []string) string) *CmdTargetBuilder {
	b.makeFunc = f
	return b
}

func (b *CmdTargetBuilder) GetService() *service.Service {
	return b.service
}

func (b *CmdTargetBuilder) GetCmd() string {
	return b.cmd
}

func (b *CmdTargetBuilder) AppendResult(nodeName, result string) {
	b.results = append(b.results, result)
}

func (b *CmdTargetBuilder) Wait(waitNum int) {
	b.waitNum = waitNum
}

func (b *CmdTargetBuilder) Done() {
	b.waitNum--

	if b.waitNum == 0 {
		b.cbFin(b.makeFunc(b.results))
	}
}

func (b *CmdTargetBuilder) Do(service *service.Service, cbFin func(string)) {
	n := app.Node
	members := n.GetCluster().GetMembers()

	b.service = service
	b.results = make([]string, 0)
	b.cbFin = cbFin

	found := false
	for k, v := range members {
		// skip master self
		if strings.Index(k, app.MasterId) > 0 {
			continue
		}

		nodeName := getNodeFromId(k)
		if nodeName != b.target {
			continue
		}

		found = true
		b.Wait(1)

		pid := actor.NewPID(fmt.Sprintf("%v:%v", v.Host, v.Port), define.NodeAdmin)
		b.doFunc(b, pid, k, v)
	}

	if !found {
		cbFin("can not find target: " + b.target)
	}
}

func DefaultTargetMakeResult(results []string) string {
	finalResult := ""
	for _, v := range results {
		finalResult += fmt.Sprintf("%v\n", v)
	}
	return finalResult
}
