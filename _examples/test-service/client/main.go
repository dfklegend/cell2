package main

import (
	"strconv"

	console "github.com/asynkron/goconsole"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"

	"github.com/dfklegend/cell2/actorex/service"
	mymsgs "test-service/messages"
	"test-service/services"
	//"test-service/messages"
)

/*
	多节点服务访问
	client->addService1.add
			sleep(1)
			-> addServer2.add
				sleep(1)
				response
			response

	敲入
	/req 1

	返回输入数字+3
*/

func main() {
	system := actor.NewActorSystem()
	config := remote.Configure("127.0.0.1", 0)
	remoter := remote.NewRemote(system, config)
	remoter.Start()

	//server := actor.NewPID("127.0.0.1:8080", "chatserver")

	// define root context
	rootContext := system.Root
	//props := actor.PropsFromProducer(func() actor.Actor { return NewClientService() },
	//	actor.WithDispatcher(service.NewAndStartDisp()))
	props, _ := service.NewServicePropsWithNewScheDisp(func() actor.Actor { return services.NewClientService() }, "")

	client := rootContext.Spawn(props)

	cons := console.NewConsole(func(text string) {
		//
	})

	// write /nick NAME to change your chat username
	cons.Command("/req", func(s string) {
		num, _ := strconv.Atoi(s)
		rootContext.Send(client, &mymsgs.Add{I: int32(num)})
	})
	cons.Run()
}
