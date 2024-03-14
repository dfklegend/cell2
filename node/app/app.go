package app

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/baseapp"
	"github.com/dfklegend/cell2/node/cluster"
	"github.com/dfklegend/cell2/node/config"
	"github.com/dfklegend/cell2/node/service"
	"github.com/dfklegend/cell2/nodectrl"
	"github.com/dfklegend/cell2/nodectrl/define"
	"github.com/dfklegend/cell2/utils/logger"
)

var (
	//	节点app
	Node *App = NewNode()
)

type App struct {
	*baseapp.App

	// 节点配置
	mainDir    string
	nodes      *config.Nodes
	nodeId     string
	nodeInfo   *config.NodeInfo
	masterInfo *config.MasterInfo

	clusterInfo *config.ClusterInfo
	cluster     *Cluster
	provider    cluster.Provider

	// actor system
	system *actor.ActorSystem

	nodeCtrl *nodectrl.NodeCtrl
}

func NewNode() *App {
	return &App{
		App:      baseapp.NewApp(),
		cluster:  NewCluster(),
		nodeCtrl: nodectrl.NewNodeCtrl(),
	}
}

func GetNode() *App {
	return Node
}

func (a *App) GetApp() *baseapp.App {
	return a.App
}

func (a *App) EnableFileLog(prefix, logDir string) {
	logger.EnableFileLog(prefix, logDir)
}

//	cfgDir: nodes.yaml路径
func (a *App) Prepare(cfgDir string) {
	logger.Log.Infof("App.Prepare %v", cfgDir)

	a.App.Prepare()
	a.nodes = config.LoadNodes(cfgDir)
	a.clusterInfo = config.LoadCluster(cfgDir)
	a.masterInfo = config.LoadMaster(cfgDir)
}

//	调用前，注册号launchmode
func (a *App) StartNode(id string, fin func(succ bool)) {
	logger.L.Infof("startnode: %v", id)

	nodes := a.nodes
	if nodes == nil {
		logger.Log.Errorf("can not find nodes")
		return
	}

	a.nodeId = id
	nodeInfo := nodes.Nodes[id]
	if nodeInfo == nil {
		logger.Log.Errorf("bad node select: %v", id)
		return
	}
	a.nodeInfo = nodeInfo
	logger.L.Infof("start node %v with cluster: %v", id, a.clusterInfo.Name)

	// init自身信息
	a.cluster.InitSelf(nodeInfo.Address, a.clusterInfo, id,
		nodeInfo.Services, nodes.Services)

	launchMode := nodeInfo.StartMode
	// 在launchmode PrepareModules中添加module
	baseapp.LaunchApp(a.App, launchMode, func(succ bool) {
		logger.Log.Infof("launch over, succ: %v\n", succ)

		// start service
		a.StartServices()
		a.StartNodeCtrl()
		fin(succ)
	})
}

func (a *App) GetNodes() *config.Nodes {
	return a.nodes
}

func (a *App) GetServiceCfg(name string) *config.ServiceInfo {
	return a.nodes.Services[name]
}

func (a *App) GetNodeId() string {
	return a.nodeId
}

func (a *App) GetNodeInfo() *config.NodeInfo {
	return a.nodeInfo
}

func (a *App) GetClusterCfg() *config.ClusterInfo {
	return a.clusterInfo
}

func (a *App) GetMasterInfo() *config.MasterInfo {
	return a.masterInfo
}

func (a *App) GetCluster() *Cluster {
	return a.cluster
}

func (a *App) GetNodeCtrl() *nodectrl.NodeCtrl {
	return a.nodeCtrl
}

func (a *App) SetProvider(provider cluster.Provider) {
	a.provider = provider
}

func (a *App) StopNode(fin func(succ bool)) {
	logger.L.Infof("StopNode")
	a.App.Stop(func(succ bool) {
		if fin != nil {
			fin(succ)
		}
	})
}

// WaitEnd 进程等待
func (a *App) WaitEnd() {
	for !a.App.GetRunService().IsStopped() {
		time.Sleep(time.Second)
	}
}

func (a *App) StartServices() {
	info := a.nodeInfo
	for i := 0; i < len(info.Services); i++ {
		name := info.Services[i]
		cfg := a.nodes.Services[name]
		if cfg == nil {
			logger.Log.Errorf("can not find service cfg: %v", name)
			continue
		}

		logger.L.Infof("start service: %v", name)
		service.Factory.Create(cfg.Type, name)
	}
}

func (a *App) FilterSelfServices(filter func(name string, cfg *config.ServiceInfo)) {
	if filter == nil {
		return
	}
	info := a.nodeInfo
	for i := 0; i < len(info.Services); i++ {
		name := info.Services[i]
		cfg := a.nodes.Services[name]
		if cfg == nil {
			continue
		}
		filter(name, cfg)
	}
}

func (a *App) StartNodeCtrl() {
	if a.clusterInfo.NodeCtrl {
		a.nodeCtrl.Start(a)
	}
	a.nodeCtrl.NodeId = a.nodeId
}

func (a *App) SetActorSystem(system *actor.ActorSystem) {
	a.system = system
}

func (a *App) GetActorSystem() *actor.ActorSystem {
	return a.system
}

func (a *App) GetService(name string) *actor.PID {
	s := a.GetCluster().GetService(name)
	if s == nil {
		return nil
	}
	return s.PID
}

func (a *App) UpdateNodeState(state int) {
	if a.provider == nil {
		return
	}
	a.provider.UpdateClusterState(state)
}

func (a *App) GetNodeState() int {
	return int(a.nodeCtrl.GetState())
}

func (a *App) IsWorkState() bool {
	return a.nodeCtrl.GetState() == define.Working
}
