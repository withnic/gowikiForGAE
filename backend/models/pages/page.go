package pages

import (
	"errors"
	"html/template"

	"net/http"

	"github.com/mjibson/goon"
)

type Page struct {
	Title    string        `datastore:"title" goon:"id"`
	Body     []byte        `datastore:"body,noindex"`
	HtmlBody template.HTML `datastore:"-"`
}

func (p *Page) Save(r *http.Request) error {
	if !validateTitle(p.Title) {
		return errors.New("Invalid title. Title is only alphabet.")
	}
	g := goon.NewGoon(r)
	e := new(Page)
	e.Title = p.Title
	e.Body = p.Body

	if _, err := g.Put(e); err != nil {
		return err
	}

	return nil
}
