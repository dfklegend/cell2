package cmds

import (
	"github.com/dfklegend/cell2/utils/convert"

	"mmo/common/cmd"
	"mmo/servers/logic_old/cmds/context"
)

type CmdEnterScene struct {
	cmd.ICmd
}

func (c *CmdEnterScene) GetName() string {
	return "enterscene"
}

func (c *CmdEnterScene) Do(ctx cmd.IContext, args []string, cb func(result string)) {
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

	sceneId := convert.TryParseInt64(args[0], 0)
	player.QueryAndEnterScene(uint64(sceneId), func(succ bool) {
		cb("succ")
	})
}
