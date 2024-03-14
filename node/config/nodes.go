package config

import (
	"github.com/spf13/viper"
)

/*
	配置节点和服务
*/

type NodeInfo struct {
	StartMode string
	Address   string
	Services  []string
}

type ServiceInfo struct {
	Type            string
	Frontend        bool
	ClientAddress   string
	WSClientAddress string
}

//	节点配置文件
type Nodes struct {
	Nodes    map[string]*NodeInfo
	Services map[string]*ServiceInfo
}

func LoadNodes(path string) *Nodes {
	v := viper.New()
	v.SetConfigName("nodes")
	v.AddConfigPath(path)
	v.SetConfigType("yaml")
	v.ReadInConfig()

	var obj Nodes
	v.Unmarshal(&obj)
	// can do some postprocess

	return &obj
}
