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
	config := remote.Configure("127.0.0.1", 1001)
	remoter := remote.NewRemote(system, config)
	remoter.Start()

	// define root context
	rootContext := system.Root

	//props := actor.PropsFromProducer(func() actor.Actor { return NewAddService2() },
	//	actor.WithDispatcher(service.NewAndStartDisp()))
	props, _ := service.NewServicePropsWithNewScheDisp(func() actor.Actor { return services.NewAddService2() }, "")

	rootContext.SpawnNamed(props, "addService2")

	cons := console.NewConsole(func(text string) {
		//
	})

	// write /nick NAME to change your chat username
	cons.Command("/nick", func(newNick string) {
		//
	})
	cons.Run()
}
