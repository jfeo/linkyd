package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"feodor.dk/linkyd/linky"
)

type JSONResponse struct {
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
	Error  string `json:"error,omitempty"`
}

type JSONBackend struct {
	templateData *template.Template
	linky        *linky.Linky
	writer       http.ResponseWriter
	request      *http.Request
	log          slog.Logger
}

func NewJSONBackend(linky *linky.Linky, w http.ResponseWriter, r *http.Request) *JSONBackend {
	return &JSONBackend{linky: linky, writer: w, request: r}
}

func (b *JSONBackend) Create() {
	var linkData struct {
		URL   string `json:"url"`
		Title string `json:"title,omitempty"`
		User  string `json:"user,omitempty"`
	}

	if data, err := io.ReadAll(b.request.Body); err != nil {
		b.writer.WriteHeader(400)
		return
	} else if err := json.Unmarshal(data, &linkData); err != nil {
		slog.Error("could not unmarshal data", "error", err)
		b.writer.WriteHeader(400)
		return
	}

	if linkData.URL == "" {
		b.writer.WriteHeader(400)
		return
	}

	if !strings.Contains(linkData.URL, "://") {
		linkData.URL = fmt.Sprintf("https://%s", linkData.URL)
	}

	if _, err := url.Parse(linkData.URL); err != nil {
		slog.Error("invalid url format", "error", err)
		b.writer.WriteHeader(400)
		return
	}

	link := b.linky.CreateLink(linkData.URL, linkData.Title, linkData.User)
	b.writeSuccess(link)
}

func (b *JSONBackend) Delete(id string) {
	link := b.linky.DeleteLink(id)
	b.writeSuccess(link)
}

func (b *JSONBackend) List() {
	b.writeSuccess(b.linky)
}

func (b *JSONBackend) As(user string) {
	b.writeSuccess(b.linky.AsUser(user))
}

func (b *JSONBackend) writeError(msg string) {
	resp := JSONResponse{
		Status: "error",
		Error:  msg,
	}

	b.writeJSON(resp)
}

func (b *JSONBackend) writeSuccess(data any) {
	resp := JSONResponse{
		Status: "success",
		Data:   data,
	}

	b.writeJSON(resp)
}

func (b *JSONBackend) writeJSON(data any) {
	if body, err := json.Marshal(data); err != nil {
		b.writer.WriteHeader(http.StatusInternalServerError)
	} else {
		b.writer.Write(body)
	}
}
