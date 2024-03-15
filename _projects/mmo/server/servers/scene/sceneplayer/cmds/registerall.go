package cmds

import (
	"mmo/common/cmd"
)

func RegisterCmds(mgr *cmd.Mgr) {
	mgr.Register(&RandSwitchLine{})
}
