package cmds

import (
	"fmt"

	"github.com/dfklegend/cell2/utils/convert"

	"mmo/common/cmd"
	"mmo/servers/logic_old/cmds/context"
)

type CmdDeleteCard struct {
	cmd.ICmd
}

func (c *CmdDeleteCard) GetName() string {
	return "deletecard"
}

func (c *CmdDeleteCard) Do(ctx cmd.IContext, args []string, cb func(result string)) {
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
	succ := player.DeleteCharCard(int32(id))
	cb(fmt.Sprintf("deletecard :%v result: %v", id, succ))
}
