package cmds

import (
	"fmt"

	"mmo/common/cmd"
	"mmo/servers/logic_old/cmds/context"
)

type CmdCreateCard struct {
	cmd.ICmd
}

func (c *CmdCreateCard) GetName() string {
	return "createcard"
}

func (c *CmdCreateCard) Do(ctx cmd.IContext, args []string, cb func(result string)) {
	if len(args) < 1 {
		cb("err args")
		return
	}

	// 请求scenem返回场景列表
	cmdCtx := ctx.(*context.CmdContext)
	player := cmdCtx.P
	if player == nil {
		cb("")
		return
	}

	name := args[0]
	id := player.CreateCharCard(name)
	cb(fmt.Sprintf("card :%v created", id))
}
