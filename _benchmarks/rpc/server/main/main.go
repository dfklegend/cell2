package main

/*
	测试与客户端的连接(cell/examples/chat-client-go)
*/

import (
	"flag"
	"fmt"
	"log"

	console "github.com/asynkron/goconsole"

	"github.com/dfklegend/cell2/baseapp"
	"github.com/dfklegend/cell2/baseapp/interfaces"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/route"
	"github.com/dfklegend/cell2/node/service"
	builder "github.com/dfklegend/cell2/nodebuilder"
	utils "github.com/dfklegend/cell2/nodeutils"
	"github.com/dfklegend/cell2/utils/common"
	"github.com/dfklegend/cell2/utils/logger"
	"server/services/chat"
	"server/services/gate"
)

func initRoutes() {
	rs := route.GetRouteService()
	rs.Register("chat", func(serviceType string, param route.IRouteParam) string {
		chatid := param.Get("chatid", "").(string)
		return chatid
	})
}

func main_old() {
	serverId := flag.String("id", "all-1", "the node id")
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

	// console
	cons := console.NewConsole(func(text string) {
	})

	common.GoPprofServe("8080")
	cons.Run()
}

func main() {
	serverId := flag.String("id", "all-1", "the node id")
	flag.Parse()

	fmt.Println("--------")
	fmt.Println("Welcome to cell2!")
	fmt.Printf("Start as %v\n", *serverId)
	fmt.Println("--------")

	builder.NewBuilder().
		ConfigLog(*serverId, "./logs").
		RegisterLaunchModes(func() {
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
		}).
		RegisterAPIs(true, func() {
			RegisterAllUserEntries()
		}).
		RegisterServiceCreators(func(factory service.IServiceFactory) {
			// 注册服务构建器
			factory.Register("gate", gate.NewCreator())
			factory.Register("chat", chat.NewCreator())
		}).
		RegisterRoutes(func(rs *route.RouteService) {
			rs.Register("chat", func(serviceType string, param route.IRouteParam) string {
				chatid := param.Get("chatid", "").(string)
				return chatid
			})
		}).
		StartApp("./testdata/config", *serverId, func(succ bool) {
			fmt.Printf("start ret: %v\n", succ)
			if succ {
				OnNodeStartSucc()
			}
		})

	// console
	cons := console.NewConsole(func(text string) {
	})

	common.GoPprofServe("8080")
	cons.Run()
}

func OnNodeStartSucc() {
}
