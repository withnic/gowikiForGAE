package pages

import (
	"errors"
	"html/template"
	"io/ioutil"
	"regexp"

	"net/http"

	"github.com/mjibson/goon"
	"github.com/russross/blackfriday"
)

type UtilFuncs interface {
	Load(title string, r *http.Request) (*Page, error)
	Parse(raw []byte) []byte
}

var Util UtilFuncs = utilFuncs{}

type utilFuncs struct{}

var validateTitleRegexp = regexp.MustCompile(`^([a-zA-Z0-9]+)$`)

func validateTitle(title string) bool {
	return validateTitleRegexp.MatchString(title)
}

func readBody(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (u utilFuncs) Parse(raw []byte) []byte {
	return blackfriday.MarkdownCommon(raw)
}

func getBodyFromDatastore(title string, r *http.Request) ([]byte, error) {
	g := goon.NewGoon(r)
	p := &Page{
		Title: title,
	}
	if err := g.Get(p); err != nil {
		return nil, err
	}
	return p.Body, nil
}

// Load returns Page struct and error
func (u utilFuncs) Load(title string, r *http.Request) (*Page, error) {
	if !validateTitle(title) {
		return nil, errors.New("Invalid title. Title is only alphabet.")
	}

	body, err := getBodyFromDatastore(title, r)

	if err != nil {
		return nil, err
	}

	output := u.Parse(body)
	return &Page{
		Title:    title,
		Body:     body,
		HtmlBody: template.HTML(output),
	}, err
}
