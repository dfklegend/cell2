package cmds

import (
	"github.com/dfklegend/cell2/node/app"
	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/waterfall"

	"mmo/common/cmd"
	mymsg "mmo/messages"
	"mmo/servers/scene/sceneplayer/cmds/context"
)

type RandSwitchLine struct {
	cmd.ICmd
}

func (c *RandSwitchLine) GetName() string {
	return "randswitchline"
}

func (c *RandSwitchLine) Do(ctx cmd.IContext, args []string, cb func(result string)) {
	// . 向scenem请求随机的线和场景
	// . 向logic发出切线请求
	cmdCtx := ctx.(*context.CmdContext)
	player := cmdCtx.Player
	if player == nil {
		cb("failed")
		return
	}

	if !player.CanSwitchLine() {
		cb("failed")
		return
	}

	ns := player.GetNodeService()
	var sceneId uint64
	waterfall.NewBuilder(ns.GetRunService().GetScheduler()).
		Next(func(callback waterfall.Callback, args ...any) {
			app.Request(ns, "scenem.remote.queryscenes", nil, &mymsg.SMQueryScenes{}, func(err error, ret any) {
				if err != nil {
					callback(true)
					return
				}
				ack := ret.(*mymsg.SMAckScenes)
				if ack == nil {
					callback(true)
					return
				}

				// rand one
				if len(ack.Scenes) == 0 {
					callback(true)
					return
				}

				sceneId = ack.Scenes[0]
				callback(false)
			})
		}).
		Next(func(callback waterfall.Callback, args ...any) {
			app.Request(ns, "scenem.remote.queryscene", nil, &mymsg.SMQueryScene{
				SceneId: sceneId,
			}, func(err error, ret any) {
				if err != nil {
					l.L.Errorf("query scene failed, err: %v", err)
					callback(true)
					return
				}
				callback(false, ret)
			})
		}).
		Next(func(callback waterfall.Callback, args ...any) {
			ack := args[0].(*mymsg.SMQuerySceneAck)
			app.Request(ns, "logic.logicremote.reqswitchline", player.GetLogicId(), &mymsg.LogicReqSwitchLine{
				UId:         player.GetId(),
				SceneServer: ack.ServiceId,
				SceneId:     sceneId,
				CfgId:       ack.CfgId,
				Token:       ack.Token,
			}, func(err error, ret any) {
				if err != nil {
					l.L.Errorf("reqswitchline failed, err: %v", err)
					callback(true)
					return
				}
				callback(false)
			})
		}).
		Final(func(err bool, args ...any) {
			cb("succ")
		}).
		Do()
}
