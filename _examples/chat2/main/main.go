package main

/*
	测试与客户端的连接(cell/examples/chat-client-go)
*/

import (
	"flag"
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	as "github.com/dfklegend/cell2/actorex/service"
	"github.com/dfklegend/cell2/baseapp"
	"github.com/dfklegend/cell2/baseapp/interfaces"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/builtin/msgs"
	"github.com/dfklegend/cell2/node/route"
	"github.com/dfklegend/cell2/node/service"
	utils "github.com/dfklegend/cell2/nodeutils"
	"github.com/dfklegend/cell2/utils/build"
	"github.com/dfklegend/cell2/utils/cmd"
	"github.com/dfklegend/cell2/utils/logger"

	"chat2/services/chat"
	"chat2/services/gate"
)

var (
	pidClient *actor.PID
)

func initRoutes() {
	rs := route.GetRouteService()
	rs.Register("chat", func(serviceType string, param route.IRouteParam) string {
		chatid := param.Get("chatid", "").(string)
		return chatid
	})
}

func main() {

	serverId := flag.String("id", "all-1", "the node id")
	flag.Parse()

	fmt.Println("--------")
	fmt.Println("Welcome to cell2!")
	fmt.Printf("Start as %v\n", *serverId)
	fmt.Println("--------")

	build.DumpInfo(Version, "", "")
	build.DumpBuildInfo()

	n := app.Node
	n.EnableFileLog(*serverId, "./logs")
	logger.Log.Infoln("")
	logger.Log.Infoln("--------------------------")
	logger.Log.Infoln("-------  app start -------")
	logger.Log.Infoln("--------------------------")
	logger.Log.Infof("Start as %v", *serverId)

	// 注册launch mode
	baseapp.LaunchFunc("allinone", func(app interfaces.IApp) {
		log.Printf("allinone mode\n")
		utils.NodeAddCommonModules(app)
		utils.NodeAddClusterModules(app)
	})
	baseapp.LaunchFunc("common", func(app interfaces.IApp) {
		log.Printf("common mode\n")
		utils.NodeAddCommonModules(app)
		utils.NodeAddClusterModules(app)
	})

	baseapp.LaunchFunc("gate", func(app interfaces.IApp) {
		log.Printf("chat mode\n")
		utils.NodeAddCommonModules(app)
		utils.NodeAddClusterModules(app)
	})

	// 注册api
	RegisterAllAPI()

	// 注册服务构建器
	service.Factory.Register("gate", gate.NewCreator())
	service.Factory.Register("chat", chat.NewCreator())

	initRoutes()

	n.Prepare("./testdata/config")
	n.StartNode(*serverId, func(succ bool) {
		fmt.Printf("start ret: %v\n", succ)
		if succ {
			OnNodeStartSucc()
		}
	})

	n.WaitEnd()
}

func OnNodeStartSucc() {
	system := app.Node.GetActorSystem()
	cmd.RegisterFuncCmd("retire", func(args []string) {
		as.DirectSendNotify(system.Root, app.Node.GetNodeCtrl().GetAdmin(),
			"ctrl.cmd", &msgs.CtrlCmd{
				Cmd: "retire",
			})
	})

	cmd.RegisterFuncCmd("exit", func(args []string) {
		as.DirectSendNotify(system.Root, app.Node.GetNodeCtrl().GetAdmin(),
			"ctrl.cmd", &msgs.CtrlCmd{
				Cmd: "exit",
			})
	})

	cmd.StartConsoleCmd()
}
