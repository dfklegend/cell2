package main

/*
	测试使用接口映射方式
*/

import (
	"flag"
	"fmt"
	"log"

	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"

	"chatapi/messages"
	"chatapi/services/chat"
	"chatapi/services/chatm"
	"chatapi/services/client"
	"chatapi/services/connector"
	"chatapi/services/gate"
	"chatapi/services/login"
	"github.com/dfklegend/cell2/baseapp"
	"github.com/dfklegend/cell2/baseapp/interfaces"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/service"
	utils "github.com/dfklegend/cell2/nodeutils"
	"github.com/dfklegend/cell2/utils/logger"
)

var (
	pidClient *actor.PID
)

func main() {

	serverId := flag.String("id", "client-0", "the node id")
	flag.Parse()

	fmt.Println("--------")
	fmt.Println("Welcome to cell2!")
	fmt.Printf("Start as %v\n", *serverId)
	fmt.Println("--------")

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
	baseapp.LaunchFunc("client", func(app interfaces.IApp) {
		log.Printf("gate mode\n")
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
	service.Factory.Register("connector", connector.NewCreator())
	service.Factory.Register("login", login.NewCreator())
	service.Factory.Register("chatm", chatmgr.NewCreator())
	service.Factory.Register("chat", chat.NewCreator())

	n.Prepare("./testdata/config")
	n.StartNode(*serverId, func(succ bool) {
		fmt.Printf("start ret: %v\n", succ)
		if succ {
			OnNodeStartSucc()
		}
	})

	// console
	cons := console.NewConsole(func(text string) {
	})

	cons.Command("/login", func(s string) {
		doLogin()
	})

	cons.Command("/say", func(s string) {
		doSay(s)
	})
	cons.Command("/nick", func(s string) {
		doNickname(s)
	})

	cons.Command("/rollbegin", func(s string) {
		doSay("/rollbegin")
	})

	cons.Command("/roll", func(s string) {
		doSay("/roll")
	})

	cons.Command("/rollend", func(s string) {
		doSay("/rollend")
	})

	cons.Run()
}

func OnNodeStartSucc() {
	// 客户端，启动client
	if app.Node.GetNodeInfo().StartMode != "client" {
		return
	}
	creator := client.NewCreator()

	name := "client"
	creator.Create(name)

	pid := actor.NewPID(app.Node.GetNodeInfo().Address, name)
	pidClient = pid
}

func doLogin() {
	if pidClient == nil {
		return
	}
	system := app.Node.GetActorSystem()

	system.Root.Send(pidClient, &messages.ClientLogin{})
}

func doSay(str string) {
	if pidClient == nil {
		return
	}
	system := app.Node.GetActorSystem()

	system.Root.Send(pidClient, &messages.ClientSay{
		Str: str,
	})
}

func doNickname(str string) {
	if pidClient == nil {
		return
	}
	system := app.Node.GetActorSystem()

	system.Root.Send(pidClient, &messages.ClientNickname{
		Name: str,
	})
}
