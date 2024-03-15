package buf

import (
	"mmo/modules/fight/script"
	"mmo/modules/fightscripts/bridge"
)

func init() {
	bridge.SetGetBufMgrFunc(func() script.IBufScriptMgr {
		return mgr
	})
}

func Visit() {
}
