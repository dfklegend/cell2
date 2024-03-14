package app

import (
	"fmt"
	"sync"

	"github.com/dfklegend/cell2/node/cluster"
	"github.com/dfklegend/cell2/node/config"
	"github.com/dfklegend/cell2/nodectrl/define"
)

//	impl ICluster
// 	提供给etcd provider使用
type Cluster struct {
	cfg     *config.ClusterInfo
	address string

	id       string
	services []string

	mutex sync.RWMutex

	clusterServices *ClusterServices
}

func NewCluster() *Cluster {
	return &Cluster{
		clusterServices: NewClusterServices(),
	}
}

func (c *Cluster) InitSelf(address string, cfg *config.ClusterInfo,
	ID string, services []string, serviceCfg map[string]*config.ServiceInfo) {
	c.address = address
	c.id = ID
	c.services = c.makeFullNameServices(services, serviceCfg)
	c.cfg = cfg
}

// 	fullName: serviceType.service
func (c *Cluster) makeFullNameServices(services []string,
	cfg map[string]*config.ServiceInfo) []string {
	result := make([]string, 0)
	for i := 0; i < len(services); i++ {
		name := services[i]
		entry := cfg[name]
		if entry == nil {
			continue
		}
		full := fmt.Sprintf("%v.%v", entry.Type, name)
		result = append(result, full)
	}
	return result
}

func (c *Cluster) GetAddress() string {
	return c.address
}

func (c *Cluster) GetName() string {
	return c.cfg.Name
}

func (c *Cluster) GetID() string {
	return c.id
}

func (c *Cluster) GetState() int {
	return Node.GetNodeState()
}

// 本地service
func (c *Cluster) GetServices() []string {
	return c.services
}

// UpdateClusterTopology 更新cluster数据
func (c *Cluster) UpdateClusterTopology(members []*cluster.Member) {
	c.clusterServices.MakeMembers(members)

	// TODO: 可以dispatch 事件
}

func (c *Cluster) GetServiceList(serviceType string) *ServiceList {
	return c.clusterServices.GetServiceList(serviceType)
}

func (c *Cluster) GetWorkServiceList(serviceType string) *ServiceList {
	return c.clusterServices.GetWorkServiceList(serviceType)
}

func (c *Cluster) GetWorkServiceNames() []string {
	return c.clusterServices.GetWorkServiceNames()
}

func (c *Cluster) GetService(name string) *ServiceItem {
	return c.clusterServices.GetService(name)
}

func (c *Cluster) GetMembers() map[string]*cluster.Member {
	return c.clusterServices.GetMembers()
}

// BuildSelfClusterTopology 用自身配置构建一个members, 使用自身cluster
func (c *Cluster) BuildSelfClusterTopology() []*cluster.Member {
	members := make([]*cluster.Member, 0)

	nodeName := fmt.Sprintf("%v@%v", c.GetName(), c.id)
	host, port, _ := splitHostPort(c.address)
	member := &cluster.Member{
		Id:       nodeName,
		Host:     host,
		Port:     int32(port),
		Services: c.services,
		State:    int(define.Working),
	}

	members = append(members, member)
	return members
}
