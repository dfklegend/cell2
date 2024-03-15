package cmds

import (
	"fmt"

	"github.com/dfklegend/cell2/node/app"

	"mmo/common/cmd"
	"mmo/messages"
	"mmo/servers/logic_old/cmds/context"
)

type CmdQueryScenes struct {
	cmd.ICmd
}

func (c *CmdQueryScenes) GetName() string {
	return "queryscenes"
}

func (c *CmdQueryScenes) Do(ctx cmd.IContext, args []string, cb func(result string)) {
	// 请求scenem返回场景列表
	cmdCtx := ctx.(*context.CmdContext)
	app.Request(cmdCtx.NS, "scenem.remote.queryscenes", nil, &messages.SMQueryScenes{}, func(err error, ret any) {
		if err != nil {
			return
		}
		ack := ret.(*messages.SMAckScenes)
		if ack == nil {
			return
		}

		result := "scenes: "
		for i := 0; i < len(ack.Scenes); i++ {
			if i > 0 {
				result += ","
			}
			result += fmt.Sprintf("%v", ack.Scenes[i])
		}
		cb(result)
	})
}
