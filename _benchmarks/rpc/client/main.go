package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/dfklegend/cell2/utils/cmd"
	"github.com/dfklegend/cell2/utils/common"
	"github.com/dfklegend/cell2/utils/logger"
	"github.com/dfklegend/cell2/utils/runservice"

	"client/client"
)

func StartCellClient(name string, detailLog bool) {
	c := client.NewChatClient(name)
	c.Start("127.0.0.1:30021")
	c.SetDetailLog(detailLog)
}

func StartCellClients(clientNum int) {
	for i := 0; i < clientNum; i++ {
		time.Sleep(100 * time.Millisecond)
		StartCellClient(fmt.Sprintf("client%v", i), false)
	}
}

/*
	连接chat2 测试tcp连接
*/
func main() {
	runservice.SetPerfLogLevel(runservice.LevelDisable)
	// 120比较极限
	var speed = flag.Float64("speed", float64(1), "the speed of send ")
	var reqNum = flag.Int("reqnum", 50000, "total send request ")
	var clientNum = flag.Int("clientnum", 100, "total client num")

	flag.Parse()

	client.SendSpeed = *speed
	client.MaxRequest = *reqNum

	//logger.SetWarnLevel()
	logger.EnableFileLog("client", "./logs")
	cmd.StartConsoleCmd()

	go StartCellClients(*clientNum)

	common.GoPprofServe("6060")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	fmt.Println("Got signal:", s)
}
