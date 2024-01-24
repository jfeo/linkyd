package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"feodor.dk/linkyd/templates"
)

type HTMXBackend struct {
	templateData *template.Template
	linky        *Linky
	writer       http.ResponseWriter
	request      *http.Request
}

func NewHTMXBackend(l *Linky, w http.ResponseWriter, r *http.Request) *HTMXBackend {
	return &HTMXBackend{
		templateData: templates.GetTemplateData(),
		linky:        l,
		writer:       w,
		request:      r,
	}
}

func (b *HTMXBackend) Create() {
	if err := b.request.ParseForm(); err != nil {
		slog.Error("invalid form data", "error", err)
		b.writer.WriteHeader(400)
		return
	}

	linkURL := b.request.Form.Get("url")
	linkTitle := b.request.Form.Get("title")

	if linkURL == "" {
		b.writer.WriteHeader(400)
		return
	}

	if !strings.Contains(linkURL, "://") {
		linkURL = fmt.Sprintf("https://%s", linkURL)
	}

	if _, err := url.Parse(linkURL); err != nil {
		slog.Error("invalid url format", "error", err)
		b.writer.WriteHeader(400)
		return
	}

	b.linky.CreateLink(linkURL, linkTitle)

	if err := b.templateData.ExecuteTemplate(b.writer, "links", b.linky); err != nil {
		b.writer.WriteHeader(http.StatusInternalServerError)
	}
}

func (b *HTMXBackend) Delete(id string) {
	b.linky.DeleteLink(id)
	if err := b.templateData.ExecuteTemplate(b.writer, "links", b.linky); err != nil {
		b.writer.WriteHeader(http.StatusInternalServerError)
	}
}

func (b *HTMXBackend) List() {
	if err := b.templateData.ExecuteTemplate(b.writer, "index", b.linky); err != nil {
		b.writer.WriteHeader(http.StatusInternalServerError)
	}
}
