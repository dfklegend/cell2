{{define "node"}}
package {{.PackageName}}

import (
	"log"

	"mlrs/b3"
	"mlrs/b3/config"
	"mlrs/b3/core"
	"mlrs/b3/factory"
)

func init() {
	factory.RegisterMap().Register("{{.NodeName}}", &{{.NodeName}}{})
}

//{{.NodeName}} {{.NodeTitle}}
type {{.NodeName}} struct {
	core.{{.Category}}
	{{range $index, $element := .FieldName}}
    {{$element}} string
    {{end}}
}

func (n *{{.NodeName}}) Initialize(params *config.BTNodeCfg) {
	n.BaseNode.Initialize(params)
	//TODO 初始化变量
}

func (n *{{.NodeName}}) OnOpen(tick *core.Tick) {
    log.Println(n.GetName())
}

func (n *{{.NodeName}}) OnTick(tick *core.Tick) b3.Status {
	return b3.SUCCESS
}

{{- end}}