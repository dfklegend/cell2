package goscript

import (
	"mmo/modules/fight/script"
)

func createGoProvicer(args ...any) script.IScriptProvider {
	return newProvider(args[0].(script.IBufScriptMgr))
}
