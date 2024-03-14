package master

import (
	"github.com/dfklegend/cell2/actorex/service"
	"github.com/dfklegend/cell2/actorex/service/servicemsgs"
	"github.com/dfklegend/cell2/node/app"
	"master/cmdservice"
	"master/consolecmd"
	"master/telnetcmd"
	"master/webservice"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) Start() {
	a.startCmdService()
	a.startConsoleCmd()
	a.startTelnetCmd()
	a.startWebService()
}

func (a *App) startCmdService() {
	system := app.Node.GetActorSystem()
	pid := cmdservice.CreateService(system)

	service.DirectSendNotify(system.Root, pid, "service.start", &servicemsgs.EmptyArg{})
}

// startConsoleCmd 控制台输入
func (a *App) startConsoleCmd() {
	system := app.Node.GetActorSystem()
	pid := consolecmd.CreateService(system)

	service.DirectSendNotify(system.Root, pid, "console.start", &servicemsgs.EmptyArg{})
}

func (a *App) startTelnetCmd() {
	system := app.Node.GetActorSystem()
	pid := telnetcmd.CreateService(system)

	service.DirectSendNotify(system.Root, pid, "console.start", &servicemsgs.EmptyArg{})
}

func (a *App) startWebService() {
	system := app.Node.GetActorSystem()
	pid := webservice.CreateService(system)

	service.DirectSendNotify(system.Root, pid, "web.start", &servicemsgs.EmptyArg{})
}
