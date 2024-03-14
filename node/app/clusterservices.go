package app

import (
	"fmt"
	"strings"

	"github.com/asynkron/protoactor-go/actor"

	"github.com/dfklegend/cell2/node/cluster"
	l "github.com/dfklegend/cell2/utils/logger"
)

type ServiceItem struct {
	// like "chat-1"
	Name string
	// clustername@nodeId
	ClusterNodeID string
	// node state, 使用方便
	State int
	// 对应PID
	PID *actor.PID
}

//	IClusterServices
//	提供对服务的查询
type IClusterServices interface {
	//	GetServiceList
	//	获取服务列表
	GetServiceList(serviceType string) *ServiceList
}

type ServiceList struct {
	Items []*ServiceItem
}

type ClusterServices struct {
	// 整体对象赋值，所以，并不存在锁的问题
	// TestSync
	//	serviceType: ServiceList
	typeServices map[string]*ServiceList
	// 实际working状态的service
	workingServices map[string]*ServiceList

	// nodeId: Member
	members map[string]*cluster.Member

	// 根据serviceName获取ServiceItem
	// 无论状态如何
	services map[string]*ServiceItem
}

func NewClusterServices() *ClusterServices {
	return &ClusterServices{
		//lock:    &sync.RWMutex{},
		typeServices: make(map[string]*ServiceList),
		members:      make(map[string]*cluster.Member),
	}
}

func (s *ClusterServices) MakeMembers(members []*cluster.Member) {
	m1, m2, m3, m4 := MakeMembers(members)

	s.members = m1
	s.typeServices = m2
	s.workingServices = m3
	s.services = m4

	s.dumpBrief()
}

func (s *ClusterServices) GetServiceList(serviceType string) *ServiceList {
	return s.typeServices[serviceType]
}

func (s *ClusterServices) GetWorkServiceList(serviceType string) *ServiceList {
	return s.workingServices[serviceType]
}

func (s *ClusterServices) GetWorkServices() []*ServiceItem {
	m := s.workingServices

	services := make([]*ServiceItem, 0)
	for _, v := range m {
		for _, v1 := range v.Items {
			services = append(services, v1)
		}
	}
	return services
}

func (s *ClusterServices) GetWorkServiceNames() []string {
	services := s.GetWorkServices()

	names := make([]string, len(services))
	for k, v := range services {
		names[k] = v.Name
	}
	return names
}

func (s *ClusterServices) GetService(name string) *ServiceItem {
	return s.services[name]
}

func (s *ClusterServices) GetMembers() map[string]*cluster.Member {
	return s.members
}

func (s *ClusterServices) dumpBrief() {
	members := s.members
	l.L.Infof("-- cluster nodes --")
	for k, _ := range members {
		l.L.Infof("   %v", k)
	}
}

//	funcs
// 	本机的Service的PID地址
func MakeMembers(members []*cluster.Member) (
	map[string]*cluster.Member,
	map[string]*ServiceList, map[string]*ServiceList,
	map[string]*ServiceItem) {

	// make new members
	// nodeId: Member
	newMembers := make(map[string]*cluster.Member)
	// serviceType: ServiceList
	newTypeList := make(map[string]*ServiceList)
	newWorkingList := make(map[string]*ServiceList)
	newServices := make(map[string]*ServiceItem)

	for i := 0; i < len(members); i++ {
		one := members[i]
		newMembers[one.Id] = one

		for j := 0; j < len(one.Services); j++ {
			addService(newTypeList, one.Id, one.State, one.Services[j])
			if IsWorkState(one.State) {
				addService(newWorkingList, one.Id, one.State, one.Services[j])
			}
		}
	}

	// 更新PID
	for _, v := range newTypeList {
		for i := 0; i < len(v.Items); i++ {
			item := v.Items[i]
			makePID(item, newMembers)

			if newServices[item.Name] != nil {
				l.L.Errorf("duplicate service name: %v", item.Name)
				continue
			}
			newServices[item.Name] = item
		}
	}

	return newMembers, newTypeList, newWorkingList, newServices
}

func addService(center map[string]*ServiceList, nodeId string, state int, fullServiceName string) {
	serviceType, name := SplitServiceName(fullServiceName)
	if serviceType == "" || name == "" {
		return
	}

	// 获取entry
	entry := center[serviceType]
	if entry == nil {
		entry = &ServiceList{
			Items: make([]*ServiceItem, 0),
		}
		center[serviceType] = entry
	}

	// 插入
	item := &ServiceItem{
		Name:          name,
		ClusterNodeID: nodeId,
		State:         state,
		PID:           nil,
	}
	entry.Items = append(entry.Items, item)
}

func makePID(item *ServiceItem, members map[string]*cluster.Member) {
	node := members[item.ClusterNodeID]
	if node == nil {
		return
	}

	address := fmt.Sprintf("%v:%v", node.Host, node.Port)
	item.PID = actor.NewPID(address, item.Name)
}

func SplitServiceName(fullName string) (string, string) {
	subs := strings.Split(fullName, ".")
	if len(subs) != 2 {
		return "", ""
	}
	return subs[0], subs[1]
}

func SplitNodeId(clusterNodeID string) (string, string) {
	subs := strings.Split(clusterNodeID, "@")
	if len(subs) != 2 {
		return "", ""
	}
	return subs[0], subs[1]
}
