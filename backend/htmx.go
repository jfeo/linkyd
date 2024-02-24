package backend

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"feodor.dk/linkyd/linky"
	"feodor.dk/linkyd/templates"
)

type HTMXBackend struct {
	templateData *template.Template
	linky        *linky.Linky
	writer       http.ResponseWriter
	request      *http.Request
}

func NewHTMXBackend(l *linky.Linky, w http.ResponseWriter, r *http.Request) *HTMXBackend {
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
	linkUser := b.request.Form.Get("user")

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

	b.linky.CreateLink(linkURL, linkTitle, linkUser)

	var templateData any
	if linkUser == "" {
		templateData = b.linky
	} else {
		templateData = b.linky.AsUser(linkUser)
	}

	slog.Debug("rendering create template", slog.Any("templateData", templateData))

	if err := b.templateData.ExecuteTemplate(b.writer, "links", templateData); err != nil {
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

func (b *HTMXBackend) As(user string) {
	if err := b.templateData.ExecuteTemplate(b.writer, "asuser", b.linky.AsUser(user)); err != nil {
		b.writer.WriteHeader(http.StatusInternalServerError)
		slog.Error("error writing asuser template", slog.Any("error", err))
	}
}
