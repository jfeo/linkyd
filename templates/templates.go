package templates

import (
	"embed"
	"log"
	"text/template"
)

//go:embed *.tmpl.html
var templatesFS embed.FS

type Templates struct {
	List   *template.Template
	AsUser *template.Template
	Links  *template.Template
}

func LoadTemplates() Templates {
	baseTmpl, err := template.New("base.tmpl.html").ParseFS(templatesFS, "base.tmpl.html")
	if err != nil {
		log.Fatal(err)
	}

	return Templates{
		List:   parseBlockTemplate(baseTmpl, "list.tmpl.html", "links.tmpl.html"),
		AsUser: parseBlockTemplate(baseTmpl, "asuser.tmpl.html", "base.tmpl.html", "links.tmpl.html"),
		Links:  template.Must(template.ParseFS(templatesFS, "links.tmpl.html")),
	}
}

func parseBlockTemplate(baseTmpl *template.Template, patterns ...string) *template.Template {
	return template.Must(template.Must(baseTmpl.Clone()).ParseFS(templatesFS, patterns...))
}
