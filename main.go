package main

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"text/template"

	"feodor.dk/linkyd/templates"
)

func getId() {

}

type Link struct {
	Id    string
	Url   string
	Title string
}

type Linky struct {
	idIncrement int
	Links       map[string]Link
}

func (l *Linky) AddLink(url string, title string) {
	log.Printf("Adding link ID=%d", l.idIncrement)
	id := strconv.Itoa(l.idIncrement)
	l.Links[id] = Link{Id: id, Url: url, Title: title}
	l.idIncrement++
	log.Printf("Next ID=%d", l.idIncrement)
}

func (l *Linky) RemoveLink(id string) {
	log.Printf("Removing link %s", id)
	delete(l.Links, id)
}

func (l Linky) RenderTemplateOr500(w http.ResponseWriter, r *http.Request, tmpl *template.Template) {
	err := tmpl.Execute(w, l)

	if err != nil {
		r.Response.StatusCode = 500
	}
}

func getPathId(r *http.Request, argumentIndex int) (string, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != argumentIndex+1 {
		return "", errors.New("Invalid path")
	}

	return parts[argumentIndex], nil
}

func main() {
	template := templates.GetTemplateData()
	linky := Linky{Links: make(map[string]Link)}
	linky.AddLink("https://google.com", "An evil search engine")
	linky.AddLink("https://duckduckgo.com", "An ethical search engine")
	linky.AddLink("https://yahoo.com", "A forgotten search engine")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template.ExecuteTemplate(w, "index", linky)
	})

	http.HandleFunc("/links", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			if r.ParseForm() != nil {
				w.WriteHeader(400)
				return
			}

			var formUrl string = r.Form.Get("url")
			_, urlErr := url.ParseRequestURI(formUrl)
			if urlErr != nil {
				w.WriteHeader(400)
				return
			}

			var formTitle string = r.Form.Get("title")
			if len(strings.Trim(formTitle, " \t\n")) == 0 {
				w.WriteHeader(400)
				return
			}

			linky.AddLink(formUrl, formTitle)
		}
		template.ExecuteTemplate(w, "links", linky)
	})

	http.HandleFunc("/links/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			w.WriteHeader(405)
			return
		}

		if id, err := getPathId(r, 2); err != nil {
			w.WriteHeader(400)
		} else {
			linky.RemoveLink(id)
		}

		template.ExecuteTemplate(w, "links", linky)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
