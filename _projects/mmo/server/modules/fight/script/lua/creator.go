package lua

import (
	"mmo/modules/fight/lua/env"
	"mmo/modules/fight/script"
)

func createLuaProvicer(args ...any) script.IScriptProvider {
	return newProvider(args[0].(*env.ScriptEnvData))
}
