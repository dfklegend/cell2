package cmds

import (
	"mmo/common/cmd"
	"mmo/servers/logic_old/cmds/context"
)

type CmdListCard struct {
	cmd.ICmd
}

func (c *CmdListCard) GetName() string {
	return "listcard"
}

func (c *CmdListCard) Do(ctx cmd.IContext, args []string, cb func(result string)) {

	// 请求scenem返回场景列表
	cmdCtx := ctx.(*context.CmdContext)
	player := cmdCtx.P
	if player == nil {
		cb("")
		return
	}

	cb(player.GetCharCardBrief())
}
