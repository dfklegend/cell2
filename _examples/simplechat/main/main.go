package main

/*
	测试使用nodeApp
*/

import (
	"flag"
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"

	nodeapi "github.com/dfklegend/cell2/apimapper/registry"
	"github.com/dfklegend/cell2/node/service"
	"simplechat/messages"
	"simplechat/services/chat"
	chatmgr "simplechat/services/chatm"
	client "simplechat/services/client"
	"simplechat/services/connector"
	"simplechat/services/gate"
	"simplechat/services/login"

	console "github.com/asynkron/goconsole"

	"github.com/dfklegend/cell2/baseapp"
	"github.com/dfklegend/cell2/baseapp/interfaces"
	"github.com/dfklegend/cell2/node/app"
	utils "github.com/dfklegend/cell2/nodeutils"
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

	// 注册服务构建器
	service.Factory.Register("gate", gate.NewCreator())
	service.Factory.Register("connector", connector.NewCreator())
	service.Factory.Register("login", login.NewCreator())
	service.Factory.Register("chatm", chatmgr.NewCreator())
	service.Factory.Register("chat", chat.NewCreator())

	utils.NodeInitSystemAPI()
	nodeapi.Registry.Build()

	n := app.Node
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
