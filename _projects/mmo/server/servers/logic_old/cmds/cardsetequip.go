package cmds

import (
	"github.com/dfklegend/cell2/utils/convert"

	"mmo/common/cmd"
	"mmo/servers/logic_old/cmds/context"
)

type CmdCardSetEquip struct {
	cmd.ICmd
}

func (c *CmdCardSetEquip) GetName() string {
	return "cardsetequip"
}

// Do index, equipId
func (c *CmdCardSetEquip) Do(ctx cmd.IContext, args []string, cb func(result string)) {
	if len(args) < 2 {
		cb("args: slot equipId")
		return
	}

	cmdCtx := ctx.(*context.CmdContext)
	player := cmdCtx.P
	if player == nil {
		cb("")
		return
	}

	slot := convert.TryParseInt64(args[0], 0)
	equipId := args[1]
	player.CardSetEquip(int(slot), equipId)
	cb("succ")
}
