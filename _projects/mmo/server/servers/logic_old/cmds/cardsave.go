package cmds

import (
	"mmo/common/cmd"
	"mmo/servers/logic_old/cmds/context"
)

type CmdSaveCard struct {
	cmd.ICmd
}

func (c *CmdSaveCard) GetName() string {
	return "savecard"
}

func (c *CmdSaveCard) Do(ctx cmd.IContext, args []string, cb func(result string)) {
	// 请求scenem返回场景列表
	cmdCtx := ctx.(*context.CmdContext)
	player := cmdCtx.P
	if player == nil {
		cb("")
		return
	}
	player.CardSave()
	cb("succ")
}
