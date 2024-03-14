package clustermodule

import (
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/dfklegend/cell2/baseapp/interfaces"
	"github.com/dfklegend/cell2/baseapp/module"
	"github.com/dfklegend/cell2/node/app"
	"github.com/dfklegend/cell2/node/cluster/clusterproviders/etcd"
	l "github.com/dfklegend/cell2/utils/logger"
)

type ClusterModule struct {
	*module.BaseModule
	provider *etcd.Provider
}

func NewClusterModule() *ClusterModule {
	return &ClusterModule{
		BaseModule: module.NewBaseModule(),
	}
}

func (c *ClusterModule) Start(next interfaces.FuncWithSucc) {
	app := app.Node
	cfg := app.GetClusterCfg()

	if !cfg.Enable {
		l.Log.Infof("    use self cluster    ")
		c.makeSelfCluster()
		next(true)
		return
	}

	provider, err := etcd.NewWithConfig("/cell2", clientv3.Config{
		Endpoints:   []string{cfg.ETCDServer},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		l.Log.Errorln("etcd NewWithConfig failed!")
		next(false)
		return
	}
	c.provider = provider
	app.SetProvider(provider)

	l.Log.Infof("start connect to ETCDServer: %v", cfg.ETCDServer)

	checking := true
	go func() {
		for checking {
			time.Sleep(5 * time.Second)
			if !checking {
				break
			}
			// 如果不需要，可以关闭此模块
			l.Log.Warnln("Please check whether ETCD server is open!")
		}
	}()

	//	如果etcd服务器未开放，会卡住
	err = provider.StartMember(app.GetCluster())
	checking = false
	if err != nil {
		l.Log.Errorln("etcd StartMember failed!")
		next(false)
	}

	l.Log.Infoln("etcd start succ!")
	next(true)
}

func (c *ClusterModule) Stop(next interfaces.FuncWithSucc) {
	if c.provider != nil {
		c.provider.Shutdown(true)
		c.provider = nil
	}

	next(true)
}

func (c *ClusterModule) makeSelfCluster() {
	cluster := app.Node.GetCluster()
	cluster.UpdateClusterTopology(cluster.BuildSelfClusterTopology())
}
