package main

import (
	"chat2/services/chat"
	"chat2/services/gate"
	nodeapi "github.com/dfklegend/cell2/apimapper/registry"
	utils "github.com/dfklegend/cell2/nodeutils"
)

func RegisterAllUserEntries() {
	gate.Visit()
	chat.Visit()
}

func RegisterAllAPI() {
	utils.NodeInitSystemAPI()
	RegisterAllUserEntries()
	nodeapi.Registry.Build()
	nodeapi.Registry.DumpAll()
}
