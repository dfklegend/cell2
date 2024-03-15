package main

/*

 */

import (
	"flag"
	"fmt"
	"log"

	as "github.com/dfklegend/cell2/actorex/service"
	"github.com/dfklegend/cell2/baseapp/interfaces"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/builtin/msgs"
	cconfig "github.com/dfklegend/cell2/node/client/impls/config"
	"github.com/dfklegend/cell2/node/route"
	"github.com/dfklegend/cell2/node/service"
	builder "github.com/dfklegend/cell2/nodebuilder"
	utils "github.com/dfklegend/cell2/nodeutils"
	"github.com/dfklegend/cell2/utils/cmd"
	"github.com/dfklegend/cell2/utils/common"
	"github.com/dfklegend/cell2/utils/golua"
	"github.com/dfklegend/cell2/utils/logger"

	"mmo/common/applog"
	"mmo/common/config"
	"mmo/modules/csv"
	"mmo/modules/fight/script/goscript"
	"mmo/modules/fight/script/lua"
	"mmo/modules/fight/skill/effect"
	"mmo/modules/fight/skill/formula"
	"mmo/modules/fightscripts/buf"
	"mmo/servers/db/dbop"
	"mmo/servers/logic_old/systems/register"
	"mmo/servers/scene/entity/ai/states"
	"mmo/servers/scene/entity/components/units"
	"mmo/servers/scene/sceneplayer/systems/scene_systems_visit"
	"mmo/servers/scene/space"
)

func main() {

	serverId := flag.String("id", "all-1", "the node id")
	flag.Parse()

	fmt.Println("--------")
	fmt.Println("Welcome to cell2!")
	fmt.Printf("Start as %v\n", *serverId)
	fmt.Println("--------")

	builder.NewBuilder().
		ConfigLog(*serverId, "./logs").
		Next(func() {
			config.InitConfig("./data/config")
			applog.InitLogs()
			service.EnableServiceLogFormatter()
		}).
		SetDefaultLaunchFunc(func(app interfaces.IApp) {
			logger.L.Infof("default mode\n")
			utils.NodeAddCommonModules(app)
			utils.NodeAddClusterModules(app)
		}).
		RegisterAPIs(true, func() {
			RegisterAllAPI()
		}).
		RegisterServiceCreators(func(factory service.IServiceFactory) {
			registerServices(factory)
		}).
		RegisterRoutes(func(rs *route.RouteService) {
			registerRoutes(rs)
		}).
		Next(func() {
			// register all "autoregister" things
			cconfig.PomeloSetProtoSerializer()
			states.Visit()
			space.Visit()
			units.Register()
			register.DoRegisterAllSystems()
			formula.RegisterAll()
			effect.RegisterAll()
			lua.Visit()
			goscript.Visit()
			buf.Visit()
			// scene systems
			scene_systems_visit.Visit()

			csv.LoadAll("data/csv")
			effect.ProcessArgs()

			LoadAllSceneCfgs("data/scenes")

			// compile all lua scripts
			golua.InitLuaPathAndCompile("./data/luascripts", true)
		}).
		OnPrepare(func() {
			dbop.SetRedisKeyPrefix(app.Node.GetClusterCfg().Name)
		}).
		StartApp("./data/config", *serverId, func(succ bool) {
			fmt.Printf("start ret: %v\n", succ)
			if succ {
				OnNodeStartSucc()
			}
		})

	common.GoPprofServe("8081")

	app.Node.WaitEnd()
}

func OnNodeStartSucc() {
	system := app.Node.GetActorSystem()

	cmd.RegisterFuncCmd("who", func(args []string) {
		log.Printf("node: %v\n", app.Node.GetNodeId())
	})

	cmd.RegisterFuncCmd("stat", func(args []string) {
		as.DirectSendNotify(system.Root, app.Node.GetNodeCtrl().GetAdmin(),
			"ctrl.cmd", &msgs.CtrlCmd{
				Cmd: "stat",
			})
	})

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
