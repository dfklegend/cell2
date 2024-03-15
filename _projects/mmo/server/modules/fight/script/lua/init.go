package lua

import (
	"mmo/modules/fight/script"
)

func init() {
	script.SetLuaProviderCreator(createLuaProvicer)
}

func Visit() {
}
