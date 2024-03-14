package app

import (
	"github.com/dfklegend/cell2/node/config"
	"github.com/dfklegend/cell2/utils/logger"
)

const (
	MasterId = "__master__"
)

// StartMaster	启动Master
func (a *App) StartMaster(cbFin func(succ bool)) {
	logger.L.Infof("start master server")
	if a.masterInfo == nil {
		return
	}

	nodes := a.nodes

	id := MasterId
	a.nodeId = id

	// make nodeInfo
	nodeInfo := &config.NodeInfo{
		Address:  a.masterInfo.Address,
		Services: make([]string, 0),
	}

	a.nodeInfo = nodeInfo

	logger.L.Infof("start master with cluster: %v", a.clusterInfo.Name)

	// init自身信息
	a.cluster.InitSelf(nodeInfo.Address, a.clusterInfo, id,
		nodeInfo.Services, nodes.Services)

	a.Start(func(succ bool) {
		// create master
		cbFin(succ)
	})
}
