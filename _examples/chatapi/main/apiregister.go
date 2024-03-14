package main

import (
	"chatapi/services/chat"
	chatmgr "chatapi/services/chatm"
	nodeapi "github.com/dfklegend/cell2/apimapper/registry"
	utils "github.com/dfklegend/cell2/nodeutils"
)

func RegisterAllUserEntries() {
	chat.Visit()
	chatmgr.Visit()
}

func RegisterAllAPI() {
	utils.NodeInitSystemAPI()
	RegisterAllUserEntries()
	nodeapi.Registry.Build()
}
