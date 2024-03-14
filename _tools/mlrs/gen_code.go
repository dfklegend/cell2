package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"mlrs/b3/config"
)

type TplData struct {
	PackageName string
	NodeName    string
	NodeTitle   string
	Category    string
	FieldName   []string
}

func genCode(tpl *template.Template, node config.BTCustomNodeCfg) {
	title := node.Title
	name := node.Name
	category := node.Category
	properties := node.Properties

	_, err := os.Stat(fmt.Sprintf("example/%v.go", name))
	if err == nil {
		return
	}

	var fields []string
	for key, _ := range properties {
		fields = append(fields, key)
	}

	file, _ := os.OpenFile(fmt.Sprintf("example/%v.go", name), os.O_CREATE|os.O_WRONLY, 0666)
	tpl.ExecuteTemplate(file, "node", TplData{
		PackageName: packageName,
		NodeName:    name,
		NodeTitle:   title,
		Category:    firstUpper(category),
		FieldName:   fields,
	})
}

func loadCodeTpl() *template.Template {
	tpl := template.Must(template.ParseGlob("tpl/*.tpl"))
	return tpl
}

func firstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
