package config

import (
	"fmt"
)

//BTNodeCfg 节点json类型
type BTNodeCfg struct {
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	Category    string                 `json:"category"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Children    []string               `json:"children"`
	Child       string                 `json:"child"`
	Parameters  map[string]interface{} `json:"parameters"`
	Properties  map[string]interface{} `json:"properties"`
}

func (this *BTNodeCfg) GetProperty(name string) float64 {
	v, ok := this.Properties[name]
	if !ok {
		panic("GetProperty err ,no vlaue:" + name)
		return 0
	}
	f64, fok := v.(float64)
	if !fok {
		fmt.Println("GetProperty err ,format not fload64:", name, v)
		panic("GetProperty err ,format not fload64:" + name)
		return 0
	}
	return f64
}

func (this *BTNodeCfg) GetPropertyAsInt(name string) int {
	v := this.GetProperty(name)
	i := int(v)
	return i
}
func (this *BTNodeCfg) GetPropertyAsInt64(name string) int64 {
	v := this.GetProperty(name)
	i := int64(v)
	return i
}
func (this *BTNodeCfg) GetPropertyAsBool(name string) bool {
	v, ok := this.Properties[name]
	if !ok {
		//panic("GetProperty err ,no vlaue:" + name)
		return false
	}

	b, fok := v.(bool)
	if !fok {
		if str, sok := v.(string); sok {
			return str == "true"
		}
		fmt.Println("GetProperty err ,format not bool:", name, v)
		panic("GetProperty err ,format not bool:" + name)
		return false
	}
	return b
}
func (this *BTNodeCfg) GetPropertyAsString(name string) string {
	v, ok := this.Properties[name]
	if !ok {
		panic("GetProperty err ,no vlaue:" + name)
		return ""
	}

	str, fok := v.(string)
	if !fok {
		fmt.Println("GetProperty err ,format not string:", name, v)
		panic("GetProperty err ,format not string:" + name)
		return ""
	}
	return str
}
