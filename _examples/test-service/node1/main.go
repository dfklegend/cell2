package main

import (
	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"

	"github.com/dfklegend/cell2/actorex/service"
	"test-service/services"
	//"test-service/messages"
)

func main() {
	system := actor.NewActorSystem()
	config := remote.Configure("127.0.0.1", 1000)
	remoter := remote.NewRemote(system, config)
	remoter.Start()

	// define root context
	rootContext := system.Root

	//props := actor.PropsFromProducer(func() actor.Actor { return NewAddService() },
	//	actor.WithDispatcher(service.NewAndStartDisp()))
	props, _ := service.NewServicePropsWithNewScheDisp(func() actor.Actor { return services.NewAddService() }, "")

	rootContext.SpawnNamed(props, "addService1")

	// node2
	actor.NewPID("127.0.0.1:1001", "addService2")

	cons := console.NewConsole(func(text string) {
		//
	})

	// write /nick NAME to change your chat username
	cons.Command("/nick", func(newNick string) {
		//
	})
	cons.Run()
}
