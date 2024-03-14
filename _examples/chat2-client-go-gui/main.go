package main

import (
	"chat2-client-go-gui/client"
)

const GatewayHost = "127.0.0.1:30021"

func main() {
	client.GUIStartup(GatewayHost)
	//client.ConsoleStartup(GatewayHost)
}
