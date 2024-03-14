package client

import (
	"chat2-client-go/protos"

	api "github.com/dfklegend/cell2/apimapper"
	client "github.com/dfklegend/cell2/pomeloclient"
	"github.com/dfklegend/cell2/utils/logger"
)

type Entry struct {
	api.APIEntry
}

func (e *Entry) OnNewUser(ctx *client.CellClient, msg *protos.OnNewUser) {
	//cclient := ctx.GetOwner().(*ChatClient)
	logger.Log.Infof("onNewUser:%v", msg.Name)
}

func (e *Entry) OnUserLeave(ctx *client.CellClient, msg *protos.OnNewUser) {
	logger.Log.Infof("OnUserLeave:%v", msg.Name)
}

func (e *Entry) OnMessage(ctx *client.CellClient, msg *protos.ChatMsg) {
	logger.Log.Infof("onMessage:   %v: %v", msg.Name, msg.Content)
}
