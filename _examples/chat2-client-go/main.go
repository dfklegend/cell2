package main

import (
	"fmt"
	"os"
	"os/signal"

	"chat2-client-go/client"
	"github.com/dfklegend/cell2/utils/build"
	"github.com/dfklegend/cell2/utils/cmd"
	"github.com/dfklegend/cell2/utils/common"
	"github.com/dfklegend/cell2/utils/logger"
)

func StartCellClient() {
	c := client.NewChatClient("client")
	c.Start("127.0.0.1:30021")
}

/*
	连接chat2 测试tcp连接
*/
func main() {
	logger.SetDebugLevel()
	build.DumpBuildInfo()
	cmd.StartConsoleCmd()

	StartCellClient()
	common.GoPprofServe("6060")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	fmt.Println("Got signal:", s)
}
