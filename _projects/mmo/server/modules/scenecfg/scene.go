package scenecfg

import (
	"github.com/spf13/viper"
)

type SceneCfg struct {
	Monsters []*Monster
	Exits    []*Exit
}

// Monster 怪物
type Monster struct {
	Id string
}

// Exit 出口
type Exit struct {
	X      float32
	Z      float32
	Radius float32
	To     int32
	ToX    float32
	ToZ    float32
}

func LoadSceneCfg(path string, name string) *SceneCfg {
	v := viper.New()
	v.SetConfigName(name)
	v.AddConfigPath(path)
	v.SetConfigType("yaml")
	v.ReadInConfig()

	var obj SceneCfg
	v.Unmarshal(&obj)

	return &obj
}
