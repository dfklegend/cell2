package main

import (
	"flag"

	console "github.com/asynkron/goconsole"
	"github.com/dfklegend/cell2/utils/common"
	"github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/runservice"

	"client/client"
)

func StartCellClient(name string) *client.ChatClient {
	c := client.NewChatClient(name)
	c.Start("127.0.0.1:30021")
	return c
}

/*
	连接chat2 测试tcp连接
*/
func main() {
	runservice.SetPerfLogLevel(runservice.LevelDisable)

	var name = flag.String("name", "client-1", "username")
	flag.Parse()

	//logger.SetWarnLevel()
	logger.EnableFileLog("client", "./logs")

	StartCellClient(*name)

	common.GoPprofServe("6060")

	// console
	cons := console.NewConsole(func(text string) {
	})

	cons.Run()
}
