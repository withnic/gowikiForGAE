package templates

import (
	"html/template"
	"log"
)

type UtilFuncs interface {
	New() *Template
}

var Util UtilFuncs = utilFuncs{}

type utilFuncs struct{}

func (u utilFuncs) New() *Template {
	t, err := template.ParseGlob(templatePattern)
	if err != nil {
		log.Fatal(err)
	}

	return &Template{
		templates: template.Must(t, nil),
	}
}
