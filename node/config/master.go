package config

import (
	"github.com/spf13/viper"
)

//	MasterInfo 节点配置
type MasterInfo struct {
	Address    string
	TelnetPort int
	HttpPort   int
}

func LoadMaster(path string) *MasterInfo {
	v := viper.New()
	v.SetConfigName("master")
	v.AddConfigPath(path)
	v.SetConfigType("yaml")
	v.ReadInConfig()

	var obj MasterInfo

	v.Unmarshal(&obj)
	// can do some postprocess

	return &obj
}
