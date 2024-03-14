package cmds

import (
	"fmt"
	"sort"
	"strings"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/actorex/service"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/cluster"
	"github.com/dfklegend/cell2/nodectrl/define"
	l "github.com/dfklegend/cell2/utils/logger"
)

type ICmdBuilder interface {
	GetService() *service.Service
	GetCmd() string
	AppendResult(nodeName string, result string)
	Done()
}

// CmdMembersBuilder 对每一个member执行，并合并结果
type CmdMembersBuilder struct {
	service  *service.Service
	doFunc   func(b ICmdBuilder, pid *actor.PID, id string, member *cluster.Member)
	cbFin    func(string)
	makeFunc func([]string) string

	waitMembers []string
	results     []string

	waitNum int

	// 缺省命令
	cmd string
}

func NewCmdMembersBuilder() *CmdMembersBuilder {
	return &CmdMembersBuilder{}
}

func (b *CmdMembersBuilder) WithDo(doFunc func(b ICmdBuilder, pid *actor.PID, id string, member *cluster.Member)) *CmdMembersBuilder {
	b.doFunc = doFunc
	return b
}

func (b *CmdMembersBuilder) WithCmdDo(cmd string) *CmdMembersBuilder {
	b.cmd = cmd
	b.doFunc = DefaultCmdDo
	b.makeFunc = MembersBuilderDefaultMakeResult
	return b
}

func (b *CmdMembersBuilder) WithMakeResult(f func(results []string) string) *CmdMembersBuilder {
	b.makeFunc = f
	return b
}
func (b *CmdMembersBuilder) GetService() *service.Service {
	return b.service
}

func (b *CmdMembersBuilder) GetCmd() string {
	return b.cmd
}

func (b *CmdMembersBuilder) AppendResult(nodeName string, result string) {
	index := -1
	for k, v := range b.waitMembers {
		if v == nodeName {
			index = k
			break
		}
	}

	if index < 0 {
		l.L.Error("")
		return
	}
	b.results[index] = result
}

func (b *CmdMembersBuilder) Wait(waitNum int) {
	b.waitNum = waitNum
}

func (b *CmdMembersBuilder) Done() {
	b.waitNum--

	if b.waitNum == 0 {
		b.cbFin(b.makeFunc(b.results))
	}
}

func (b *CmdMembersBuilder) Do(service *service.Service, cbFin func(string)) {
	n := app.Node
	members := n.GetCluster().GetMembers()

	b.service = service

	b.cbFin = cbFin

	waitNum := len(members) - 1

	if waitNum == 0 {
		cbFin("no members")
		return
	}

	b.Wait(waitNum)

	b.waitMembers = make([]string, 0)

	for k, _ := range members {
		// skip master self
		if strings.Index(k, app.MasterId) > 0 {
			continue
		}
		b.waitMembers = append(b.waitMembers, getNodeFromId(k))
	}

	// sort

	waitMembers := b.waitMembers
	sort.Slice(waitMembers, func(a, b int) bool {
		return strings.Compare(waitMembers[a], waitMembers[b]) < 0
	})

	b.results = make([]string, waitNum)

	for k, v := range members {
		// skip master self
		if strings.Index(k, app.MasterId) > 0 {
			continue
		}

		func(id string, v *cluster.Member) {
			pid := actor.NewPID(fmt.Sprintf("%v:%v", v.Host, v.Port), define.NodeAdmin)
			b.doFunc(b, pid, id, v)
		}(k, v)
	}
}

// MembersBuilderDefaultMakeResult 组织在一起
func MembersBuilderDefaultMakeResult(results []string) string {
	finalResult := ""
	for _, v := range results {
		finalResult += fmt.Sprintf("%v\n", v)
	}
	return finalResult
}
