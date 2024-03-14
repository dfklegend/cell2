package config

//BTCustomNodeCfg 自定义节点信息
type BTCustomNodeCfg struct {
	Version    string                 `json:"version"`
	Scope      string                 `json:"scope"`
	Name       string                 `json:"name"`
	Category   string                 `json:"category"`
	Title      string                 `json:"title"`
	Desc       string                 `json:"description"`
	Properties map[string]interface{} `json:"properties"`
	Parent     string                 `json:"parent"`
}
