package config

import (
	"github.com/dfklegend/cell2/utils/serialize"
	"github.com/dfklegend/cell2/utils/serialize/json"
	"github.com/dfklegend/cell2/utils/serialize/proto"
)

var (
	theConfig = &ClientConfig{}
)

// 定义客户端的序列化器
func init() {
	theConfig.WithSerializer(json.GetDefaultSerializer())
}

type ClientConfig struct {
	Serializer serialize.Serializer
}

func GetConfig() *ClientConfig {
	return theConfig
}

func (c *ClientConfig) WithSerializer(serializer serialize.Serializer) *ClientConfig {
	c.Serializer = serializer
	return c
}

func PomeloSetSerializer(serializer serialize.Serializer) {
	GetConfig().WithSerializer(serializer)
}

func PomeloSetProtoSerializer() {
	PomeloSetSerializer(proto.GetDefaultSerializer())
}
