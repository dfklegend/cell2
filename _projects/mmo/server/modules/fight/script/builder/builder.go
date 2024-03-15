package builder

import (
	"mmo/modules/fight/script"
	"mmo/modules/fight/script/mgr"
)

func CreateScriptMgr() script.IScriptMgr {
	return mgr.NewMgr()
}
