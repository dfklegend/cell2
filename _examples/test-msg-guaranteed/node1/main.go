package main

import (
	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"

	"github.com/dfklegend/cell2/actorex/service"
	l "github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/logger/proxy"
	"test-service/services"
	//"test-service/messages"
)

func main() {

	proxy.GetLogs().EnableFileLog("node1", "./logs")

	l.L.Infof("node1 start")

	system := actor.NewActorSystem()
	config := remote.Configure("127.0.0.1", 1000)
	remoter := remote.NewRemote(system, config)
	remoter.Start()

	// define root context
	rootContext := system.Root

	props, _ := service.NewServicePropsWithNewScheDisp(func() actor.Actor { return services.NewAddService() }, "")

	rootContext.SpawnNamed(props, "addService1")

	cons := console.NewConsole(func(text string) {
		//
	})

	// write /nick NAME to change your chat username
	cons.Command("/nick", func(newNick string) {
		//
	})
	cons.Run()
}
