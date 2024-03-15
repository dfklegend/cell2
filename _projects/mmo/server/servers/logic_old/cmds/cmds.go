package cmds

import (
	"mmo/common/cmd"
)

func RegisterCmds(mgr *cmd.Mgr) {
	mgr.Register(cmd.NewFuncCmd("test", func(ctx cmd.IContext, args []string, cb func(result string)) {
		cb("test ack")
	}))

	mgr.Register(&CmdQueryScenes{})
	mgr.Register(&CmdEnterScene{})
	mgr.Register(&CmdCreateCard{})
	mgr.Register(&CmdListCard{})
	mgr.Register(&CmdDeleteCard{})
	mgr.Register(&CmdOpenCard{})
	mgr.Register(&CmdCardSetEquip{})
	mgr.Register(&CmdSaveCard{})
}
