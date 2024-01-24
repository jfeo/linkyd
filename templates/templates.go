package templates

import (
	"embed"
	"text/template"
)

//go:embed *.tmpl.html
var templatesFS embed.FS

func GetTemplateData() *template.Template {
	base := template.New("base")
	template.Must(base.ParseFS(templatesFS, "*.tmpl.html"))
	return base
}
