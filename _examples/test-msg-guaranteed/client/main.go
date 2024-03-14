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


	敲入
	/req 1

	测试如果目标服务器暂时不可用，消息是否会丢失
	结果是，消息会在连接上后，继续发送
	(最大限度的会重发，当然内存会有消耗)
	(极端情况下，丢失，接收者收到后，进程crash了，会丢失
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
