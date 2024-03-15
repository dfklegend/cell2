package cmds

import (
	"github.com/dfklegend/cell2/utils/convert"

	"mmo/common/cmd"
	"mmo/servers/logic_old/cmds/context"
)

type CmdOpenCard struct {
	cmd.ICmd
}

func (c *CmdOpenCard) GetName() string {
	return "opencard"
}

func (c *CmdOpenCard) Do(ctx cmd.IContext, args []string, cb func(result string)) {
	if len(args) < 1 {
		cb("err args")
		return
	}

	cmdCtx := ctx.(*context.CmdContext)
	player := cmdCtx.P
	if player == nil {
		cb("")
		return
	}

	id := convert.TryParseInt64(args[0], 0)
	player.OpenCharCard(int32(id))
	cb("succ")
}
