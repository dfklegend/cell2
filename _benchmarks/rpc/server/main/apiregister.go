package main

import (
	nodeapi "github.com/dfklegend/cell2/apimapper/registry"
	utils "github.com/dfklegend/cell2/nodeutils"
	"server/services/chat"
	"server/services/gate"
)

func RegisterAllUserEntries() {
	gate.Visit()
	chat.Visit()
}

func RegisterAllAPI() {
	utils.NodeInitSystemAPI()
	RegisterAllUserEntries()
	nodeapi.Registry.Build()
}
