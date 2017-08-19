package backend

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"models/pages"
	"models/templates"
)

var renderer *templates.Template

func init() {
	renderer = templates.Util.New()
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/api/mkd", markdownHandler)
	http.HandleFunc("/", handler)
}

var validPath = regexp.MustCompile(`^/(edit|save|view)/([a-zA-Z0-9]+)$`)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
	} else {
		errorHandler(w, r, http.StatusNotFound)
	}
	return
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		err := renderer.Render(w, "errorPage", "404 Page Not found.")
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if status == http.StatusInternalServerError {
		err := renderer.Render(w, "errorPage", "Internal Server Error")
		if err != nil {
			log.Fatal(err)
		}
		return
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		fn(w, r, m[2])
	}
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := pages.Util.Load(title, r)

	if err != nil {
		p = &pages.Page{
			Title: title,
		}
	}

	err = renderer.Render(w, "editPage", p)
	if err != nil {
		log.Fatal(err)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := pages.Util.Load(title, r)

	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderer.Render(w, "viewPage", p)
}

func markdownHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	body := r.FormValue("body")
	output := pages.Util.Parse([]byte(body))

	if err := json.NewEncoder(w).Encode(string(output)); err != nil {
		log.Fatal(err)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")

	p := &pages.Page{
		Title: title,
		Body:  []byte(body),
	}

	err := p.Save(r)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
