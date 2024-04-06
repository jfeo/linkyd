package backend

import (
	"errors"
	"net/http"

	"feodor.dk/linkyd/linky"
)

var ErrUnsupportedContentType = errors.New("unsupported content type")

type LinkyBackend interface {
	Create()
	List()
	Delete(id string)
	As(user string)
}

func Get(l *linky.LinkService, w http.ResponseWriter, r *http.Request) LinkyBackend {
	switch r.Header.Get("Content-Type") {
	case "application/json":
		return NewJSONBackend(l, w, r)
	default:
		return NewHTMXBackend(l, w, r)
	}
}
