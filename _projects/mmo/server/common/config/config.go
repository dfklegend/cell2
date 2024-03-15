package config

import (
	"github.com/spf13/viper"
)

var (
	Cfg *Info
)

// 每个服务器类型配置

type GateInfo struct {
}

type ChatInfo struct {
}

type DBInfo struct {
	Redis string
}

type Info struct {
	DB   *DBInfo
	Gate *GateInfo
	Chat *ChatInfo
}

func LoadConfig(path string) *Info {
	v := viper.New()
	v.SetConfigName("cfg")
	v.AddConfigPath(path)
	v.SetConfigType("yaml")
	v.ReadInConfig()

	var obj Info

	v.Unmarshal(&obj)
	// can do some postprocess

	return &obj
}

func InitConfig(path string) {
	Cfg = LoadConfig(path)
}
