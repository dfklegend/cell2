package main

/*
	master 服务器
	显示节点状态
	可以向指定node发送命令，并显示结果

	cmds
	stat
		显示各节点状态
	nodes
		显示每个节点的服务
	services
		显示所有服务
	retire nodeid
		要求节点退休
	exit nodeid
		要求节点退出(处于retired状态的节点才可以退休)

*/

import (
	"flag"

	"github.com/dfklegend/cell2/node/app"
	builder "github.com/dfklegend/cell2/nodebuilder"
	master "master/app"
)

func main() {
	flag.Parse()

	n := app.Node

	builder.NewMasterBuilder().
		ConfigLog("master", "./logs").
		RegisterAPIs(true, func() {
			master.RegisterAllAPIEntries()
		}).
		StartMaster("./data/config", func(succ bool) {
			OnNodeStartSucc()
		})

	n.WaitEnd()
}

func OnNodeStartSucc() {
	master.NewApp().Start()
}
