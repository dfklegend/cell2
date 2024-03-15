package bridge

import (
	"mmo/modules/fight/script"
)

type FnGetBufMgr func() script.IBufScriptMgr

var fnGetBufMgr FnGetBufMgr

func SetGetBufMgrFunc(fn FnGetBufMgr) {
	fnGetBufMgr = fn
}

func GetBufMgr() script.IBufScriptMgr {
	return fnGetBufMgr()
}
