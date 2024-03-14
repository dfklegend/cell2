package cmds

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dfklegend/cell2/actorex/service"
	"github.com/dfklegend/cell2/node/app"
	"master/cmdservice/cmd"
)

func InitCmds(mgr *cmd.Mgr) {
	// 所有member广播命令
	mgr.Register("stat", doStat)

	// 对单节点操作命令
	mgr.Register("retire", doRetire)
	mgr.Register("exit", doExit)

	// 单独统计
	// 列出当前有效services
	mgr.Register("nodes", doNodes)
	mgr.Register("services", doServices)
}

func doStat(ctx cmd.IContext, args []string, cb func(string)) {
	NewCmdMembersBuilder().
		WithCmdDo("stat").
		Do(ctx.(*service.Service), cb)
}

func doServices(ctx cmd.IContext, args []string, cb func(string)) {
	n := app.Node
	names := n.GetCluster().GetWorkServiceNames()

	sort.Slice(names, func(a, b int) bool {
		return strings.Compare(names[a], names[b]) < 0
	})

	var builder strings.Builder
	for _, v := range names {
		builder.WriteString(fmt.Sprintf("  %v\n", v))
	}

	cb(builder.String())
}

func doCmdToTarget(ctx cmd.IContext, cmd string, args []string, cb func(string)) {
	if len(args) < 1 {
		cb(fmt.Sprintf("usage: %v nodeid", cmd))
		return
	}
	NewCmdTargetBuilder().
		WithCmd(cmd, args[0]).
		Do(ctx.(*service.Service), cb)
}

func doRetire(ctx cmd.IContext, args []string, cb func(string)) {
	doCmdToTarget(ctx, "retire", args, cb)
}

func doExit(ctx cmd.IContext, args []string, cb func(string)) {
	doCmdToTarget(ctx, "exit", args, cb)
}

// 列出node所属service
func doNodes(ctx cmd.IContext, args []string, cb func(string)) {
	n := app.Node

	members := n.GetCluster().GetMembers()

	names := make([]string, 0)
	for k, _ := range members {
		names = append(names, k)
	}

	sort.Slice(names, func(a, b int) bool {
		return strings.Compare(names[a], names[b]) < 0
	})

	var builder strings.Builder
	for _, name := range names {
		if strings.Index(name, app.MasterId) > 0 {
			continue
		}

		_, nodeName := app.SplitNodeId(name)
		builder.WriteString(fmt.Sprintf("%v:\n", nodeName))
		member := members[name]
		if member == nil {
			continue
		}

		for _, service := range member.Services {
			_, serviceName := app.SplitServiceName(service)
			builder.WriteString(fmt.Sprintf("  %v\n", serviceName))
		}
	}

	cb(builder.String())
}
