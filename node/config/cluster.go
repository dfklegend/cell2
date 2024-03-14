package config

import (
	"github.com/spf13/viper"
)

//	ClusterInfo cluster配置文件
type ClusterInfo struct {
	Enable     bool // false, 则不连接ectd,使用当前节点的service
	NodeCtrl   bool // 是否需要节点控制
	Name       string
	ETCDServer string

	Token string
}

func LoadCluster(path string) *ClusterInfo {
	v := viper.New()
	v.SetConfigName("cluster")
	v.AddConfigPath(path)
	v.SetConfigType("yaml")
	v.ReadInConfig()

	var obj *ClusterInfo = &ClusterInfo{
		NodeCtrl: true,
	}

	v.Unmarshal(obj)
	// can do some postprocess

	return obj
}
