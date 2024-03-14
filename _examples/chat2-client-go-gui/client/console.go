package client

import (
	"os"
	"time"

	console "github.com/asynkron/goconsole"
)

//ConsoleStartup 控制台启动
func ConsoleStartup(host string) {

	//连接 gateway
	start(host)

	for !cfg.Login {
		time.Sleep(10)
	}
	client := cfg.Client

	newConsole := console.NewConsole(client.sendChatMsg)

	newConsole.Command("exit", func(cmd string) {
		client.PomeloClient.Close()
		os.Exit(0)
	})
	newConsole.Run()
}
