package templates

import (
	"html/template"
	"net/http"
)

const (
	templatePattern = "tmpl/*.html"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w http.ResponseWriter, tmpl string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, tmpl, data)
}
