package config

//RawProjectCfg b3文件格式
type RawProjectCfg struct {
	Name string       `json:"name"`
	Data BTProjectCfg `json:"data"`
	Path string       `json:"path"`
	Desc string       `json:"description"`
}
