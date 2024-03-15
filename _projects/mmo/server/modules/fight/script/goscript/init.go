package goscript

import (
	"mmo/modules/fight/script"
)

func init() {
	script.SetGoProviderCreator(createGoProvicer)
}

func Visit() {
}
