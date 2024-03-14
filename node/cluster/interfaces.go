package cluster

//	当前服务发现问题:
//		要求本地应该有所有service的配置
//		改进:	向服务器注册时，应该将该service附带的配置也序列化过去
//				提供一个service库，用于查询
//
//				另一个思路是，service自己向管理的service同步自己的信息
//				更符合负载管理的目标
//				选用方案2

type Provider interface {
	StartMember(cluster ICluster) error
	StartClient(cluster ICluster) error
	Shutdown(graceful bool) error

	// UpdateClusterState 状态变化了
	UpdateClusterState(state int) error
}

type ICluster interface {
	//	actor remote的地址
	GetAddress() string
	// 	cluster名字
	GetName() string
	//	区分节点的Id
	GetID() string
	GetState() int

	GetServices() []string

	UpdateClusterTopology([]*Member)
}

type Member struct {
	// clusterName@nodeId
	Id   string
	Host string
	Port int32
	//	服务列表
	Services []string

	// 节点状态
	State int
}

// INodeData 节点数据接口
type INodeData interface {
	GetNodeId() string

	GetNodeData(key string) string

	SetNodeData(key string, value string)
}
