package main

/*
	测试使用nodeApp
*/

import (
	"flag"
	"fmt"
	"log"

	console "github.com/asynkron/goconsole"

	"github.com/dfklegend/cell2/baseapp"
	"github.com/dfklegend/cell2/baseapp/interfaces"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/service"
	utils "github.com/dfklegend/cell2/nodeutils"

	"test-cluster/services"
)

func main() {

	serverId := flag.String("id", "node-1", "the node id")
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
	baseapp.LaunchFunc("gate", func(app interfaces.IApp) {
		log.Printf("gate mode\n")
		utils.NodeAddCommonModules(app)
		utils.NodeAddClusterModules(app)
	})

	baseapp.LaunchFunc("chat", func(app interfaces.IApp) {
		log.Printf("chat mode\n")
		utils.NodeAddCommonModules(app)
		utils.NodeAddClusterModules(app)
	})

	// 注册服务构建器
	service.Factory.Register("gate", services.NewGateCreator())
	service.Factory.Register("chat", services.NewChatCreator())

	n := app.Node
	n.Prepare("./testdata/config")
	n.StartNode(*serverId, func(succ bool) {
		fmt.Printf("start ret: %v\n", succ)
	})

	// console
	cons := console.NewConsole(func(text string) {
	})

	cons.Command("/req", func(s string) {
	})
	cons.Run()
}
