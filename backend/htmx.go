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
	templates templates.Templates
	linky     *linky.LinkService
	writer    http.ResponseWriter
	request   *http.Request
}

func NewHTMXBackend(l *linky.LinkService, w http.ResponseWriter, r *http.Request) *HTMXBackend {
	return &HTMXBackend{
		templates: templates.LoadTemplates(),
		linky:     l,
		writer:    w,
		request:   r,
	}
}

func (b *HTMXBackend) Create() {
	slog.Info("htmx create")

	if err := b.request.ParseForm(); err != nil {
		slog.Error("invalid form data", "error", err)
		b.writer.WriteHeader(400)
		return
	}

	linkURL := b.request.Form.Get("url")
	linkTitle := b.request.Form.Get("title")
	linkUser := b.request.Form.Get("user")
	slog.Debug("htmx create got link data", slog.String("linkURL", linkURL), slog.String("linkUser", linkUser))

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

	var data any
	if linkUser == "" {
		data = b.linky.All()
	} else {
		data = b.linky.AsUser(linkUser)
	}

	if err := b.executeTemplate(b.templates.Links, data); err != nil {
		b.writer.WriteHeader(http.StatusInternalServerError)
	}
}

func (b *HTMXBackend) Delete(id string) {
	slog.Info("htmx create")

	b.linky.DeleteLink(id)
	if err := b.executeTemplate(b.templates.Links, b.linky.All()); err != nil {
		b.writer.WriteHeader(http.StatusInternalServerError)
	}
}

func (b *HTMXBackend) List() {
	slog.Info("htmx list")

	if err := b.executeTemplate(b.templates.List, b.linky.All()); err != nil {
		b.writer.WriteHeader(http.StatusInternalServerError)
	}
}

func (b *HTMXBackend) As(user string) {
	slog.Info("htmx as")

	if err := b.executeTemplate(b.templates.AsUser, b.linky.AsUser(user)); err != nil {
		b.writer.WriteHeader(http.StatusInternalServerError)
	}
}

func (b *HTMXBackend) executeTemplate(template *template.Template, data any) error {
	err := template.Execute(b.writer, data)
	name := template.Name()

	slog.Debug(
		"executing template",
		slog.String("name", name),
		slog.Any("data", data),
	)

	if err != nil {
		slog.Error(
			"error executing template",
			slog.String("error", err.Error()),
			slog.String("name", name),
			slog.Any("data", data),
		)
	}

	return err
}
